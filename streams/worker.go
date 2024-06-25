package streams

import (
	"context"
	errorshs "github.com/enviodev/hypersync-client-go/errors"
	"github.com/enviodev/hypersync-client-go/logger"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"math/big"
	"sync"
)

// WorkerFn defines a generic function type that takes a descriptor of type T and returns
// a response of type R and an error.
type WorkerFn[T any, R *types.QueryResponse] func(descriptor T) (R, error)

// Worker represents a type-specific worker that processes descriptors using a provided
// WorkerFn. It manages the processing state and results.
type Worker[T any, R *types.QueryResponse] struct {
	ctx      context.Context
	opts     *options.StreamOptions
	iterator *BlockIterator
	done     chan struct{}
	result   chan OrderedResult[T, R]
	channel  chan R
	ackCh    chan struct{} // Acknowledgment channel
	wg       sync.WaitGroup
}

// OrderedResult holds the result of processing a descriptor, including its index, the
// response, and any error encountered.
type OrderedResult[T any, R *types.QueryResponse] struct {
	index  int
	record R
	err    error
}

// NewWorker creates a new instance of a type-specific Worker.
func NewWorker[T any, R *types.QueryResponse](ctx context.Context, iterator *BlockIterator, channel chan R, done chan struct{}, opts *options.StreamOptions) (*Worker[T, R], error) {
	return &Worker[T, R]{
		ctx:      ctx,
		opts:     opts,
		iterator: iterator,
		channel:  channel,
		done:     done,
		result:   make(chan OrderedResult[T, R], big.NewInt(0).Mul(opts.Concurrency, big.NewInt(10)).Uint64()),
		ackCh:    make(chan struct{}, big.NewInt(0).Mul(opts.Concurrency, big.NewInt(10)).Uint64()), // Buffered channel for acknowledgments
	}, nil
}

// Start begins the worker's operation using the provided WorkerFn and a channel of descriptors.
func (w *Worker[T, R]) Start(workerFn WorkerFn[T, R], descriptor <-chan T) error {
	g, ctx := errgroup.WithContext(w.ctx)

	// Create an indexed channel to preserve order
	type indexedDescriptor struct {
		index int
		value T
	}
	indexedChan := make(chan indexedDescriptor)

	// Goroutine to index descriptors
	go func() {
		index := 0
		for entry := range descriptor {
			indexedChan <- indexedDescriptor{index: index, value: entry}
			index++
		}
		close(indexedChan)
	}()

	// Start worker goroutines
	for workerId := uint64(0); workerId < w.opts.Concurrency.Uint64(); workerId++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-w.done:
					return nil
				case entry, ok := <-indexedChan:
					if !ok {
						return nil
					}

					w.wg.Add(1)
					resp, err := workerFn(entry.value)
					w.result <- OrderedResult[T, R]{index: entry.index, record: resp, err: err}
				}
			}
		})
	}

	ackCount := 0 // Acknowledgment counter

	// Collect results in order and publish them to the output channel
	g.Go(func() error {
		results := make(map[int]R)
		nextIndex := 0
	mainLoop:
		for res := range w.result {
			if res.err != nil {
				logger.L().Error(
					"error processing stream entry",
					zap.Error(res.err),
					zap.Any("processing_index", res.index),
				)
				w.wg.Done()
				continue
			}
			results[res.index] = res.record
			ackCount++

			// Push results to the output channel in order
			for {
				if record, ok := results[nextIndex]; ok {
					w.channel <- record
					delete(results, nextIndex)

					// Check if this is the last record to process
					if (*record).NextBlock.Cmp(w.iterator.GetEndAsBigInt()) == 0 {
						w.wg.Done()
						break mainLoop
					}

					nextIndex++
					w.wg.Done()
				} else {
					break
				}
			}
		}

		// Wait for all messages to be processed
		w.wg.Wait()

		// Close the result channel to signal completion
		close(w.result)
		return errorshs.ErrWorkerCompleted
	})

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil && !errors.Is(err, errorshs.ErrWorkerCompleted) {
		return err
	}

	// Wait for all acknowledgments
	if !w.opts.DisableAcknowledgements {
		for i := 0; i < ackCount; i++ {
			<-w.ackCh
		}
	}

	return w.Stop()
}

// Ack acknowledges that a response has been processed.
func (w *Worker[T, R]) Ack() {
	w.ackCh <- struct{}{}
}

// Done returns a channel that can be used to signal when the worker's operations are done.
func (w *Worker[T, R]) Done() <-chan struct{} {
	return w.done
}

// Stop stops the worker's operations and waits for all goroutines to complete.
func (w *Worker[T, R]) Stop() error {
	w.wg.Wait()
	close(w.done)
	return nil
}
