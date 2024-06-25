package hypersyncgo

import (
	"context"
	"math/big"
	"sync"

	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/streams"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Stream represents a streaming process that handles data queries and responses
// using a client and worker for concurrent processing.
type Stream struct {
	ctx      context.Context
	cancelFn context.CancelFunc
	client   *Client
	queryCh  chan *types.Query
	ch       chan *types.QueryResponse
	errCh    chan error
	opts     *options.StreamOptions
	query    *types.Query
	iterator *streams.BlockIterator
	worker   *streams.Worker[*types.Query, *types.QueryResponse]
	done     chan struct{}
	mu       *sync.RWMutex
	nextIdx  uint64
	step     uint64
}

// NewStream creates a new Stream instance with the provided context, client, query, and options.
func NewStream(ctx context.Context, client *Client, query *types.Query, opts *options.StreamOptions) (*Stream, error) {
	if vErr := opts.Validate(); vErr != nil {
		return nil, errors.Wrap(vErr, "failed to validate stream options")
	}

	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	ch := make(chan *types.QueryResponse, opts.Concurrency.Uint64())
	step := opts.BatchSize.Uint64() | (uint64(0) << 32)
	blockIter := streams.NewBlockIterator(query.FromBlock.Uint64(), query.ToBlock.Uint64(), &step)
	worker, err := streams.NewWorker[*types.Query, *types.QueryResponse](ctx, blockIter, ch, done, opts)
	if err != nil {
		cancel()
		return nil, errors.Wrap(err, "failed to create new stream subscriber worker")
	}

	return &Stream{
		ctx:      ctx,
		opts:     opts,
		client:   client,
		cancelFn: cancel,
		query:    query,
		iterator: blockIter,
		worker:   worker,
		queryCh:  make(chan *types.Query, opts.Concurrency.Uint64()),
		ch:       ch,
		errCh:    make(chan error, opts.Concurrency.Uint64()),
		done:     done,
		mu:       &sync.RWMutex{},
		step:     step,
	}, nil
}

// ProcessNextQuery processes the next query using the client and returns the response or error.
func (s *Stream) ProcessNextQuery(query *types.Query) (*types.QueryResponse, error) {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	return s.client.GetArrow(ctx, query)
}

// Subscribe starts the streaming process, initializing the first query and handling subsequent ones.
func (s *Stream) Subscribe() error {
	g, ctx := errgroup.WithContext(s.ctx)

	// Initial fetch to get the first block and with it next paginated starting position
	response, err := s.client.GetArrow(ctx, s.query)
	if err != nil {
		return err
	}
	s.ch <- response

	// We've fetched everything that's requested. Considering this stream as completed.
	if response.NextBlock.Cmp(s.query.ToBlock) == 0 {
		close(s.done)
		return nil
	}

	// Start the worker to fetch remaining pages
	g.Go(func() error {
		return s.worker.Start(s.ProcessNextQuery, s.queryCh)
	})

	go func() {
		for {
			start, end, ok := s.iterator.Next()
			if !ok {
				break
			}

			iQuery := *s.query
			iQuery.FromBlock = new(big.Int).SetUint64(start)
			iQuery.ToBlock = new(big.Int).SetUint64(end)
			s.queryCh <- &iQuery

			// Exit this routine just in case at this point...
			if s.iterator.Completed() {
				return
			}
		}
	}()

	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return s.ctx.Err()
	default:
		if wErr := g.Wait(); wErr != nil {
			return wErr
		}
		return nil
	}
}

// Unsubscribe stops the stream and closes all channels associated with it.
func (s *Stream) Unsubscribe() error {
	s.worker.Stop()
	close(s.queryCh)
	close(s.ch)
	close(s.errCh)
	s.cancelFn()
	return nil
}

// QueueError adds an error to the stream's error channel.
func (s *Stream) QueueError(err error) {
	s.errCh <- err
}

// Err returns the stream's error channel.
func (s *Stream) Err() <-chan error {
	return s.errCh
}

// Channel returns the stream's response channel.
func (s *Stream) Channel() <-chan *types.QueryResponse {
	return s.ch
}

// Ack acknowledges that a response has been processed.
// This method is thread-safe.
func (s *Stream) Ack() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.worker.Ack()
}

// Done returns a channel that signals when the stream is done.
func (s *Stream) Done() <-chan struct{} {
	return s.worker.Done()
}
