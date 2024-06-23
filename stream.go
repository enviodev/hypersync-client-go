package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"golang.org/x/sync/errgroup"
)

type Stream struct {
	ctx      context.Context
	cancelFn context.CancelFunc
	client   *Client
	ch       chan *types.QueryResponse
	errCh    chan error
	opts     *options.StreamOptions
	query    *types.Query
	done     chan struct{}
}

func NewStream(ctx context.Context, client *Client, query *types.Query, opts *options.StreamOptions) (*Stream, error) {
	ctx, cancel := context.WithCancel(ctx)
	return &Stream{
		ctx:      ctx,
		client:   client,
		cancelFn: cancel,
		query:    query,
		ch:       make(chan *types.QueryResponse, opts.Concurrency.Uint64()),
		errCh:    make(chan error, opts.Concurrency.Uint64()),
		done:     make(chan struct{}),
	}, nil
}

func (s *Stream) Subscribe() error {
	g, ctx := errgroup.WithContext(s.ctx)

	g.Go(func() error {
		// Process initial response - we need to know from where to continue...
		iCtx, iCancel := context.WithCancel(s.ctx)
		defer iCancel()

		response, err := s.client.GetArrow(iCtx, s.query)
		if err != nil {
			return err
		}

		// Push initial response to the channel
		s.ch <- response

		// If we've reached a total count of requested blocks, stop and signal completion.
		if response.NextBlock.Cmp(s.query.ToBlock) == 0 {
			close(s.done)
			return nil
		}

		nextBlock := response.NextBlock
		for response.HasNextBlock() {
			iQuery := s.query
			iQuery.FromBlock = nextBlock

			iResponse, iErr := s.client.GetArrow(iCtx, iQuery)
			if iErr != nil {
				return iErr
			}

			s.ch <- iResponse

			// If we've reached a total count of requested blocks, stop and signal completion.
			if iResponse.NextBlock.Cmp(s.query.ToBlock) == 0 {
				close(s.done)
				return nil
			}

			nextBlock = iResponse.NextBlock // Update nextBlock for the next iteration
			response = iResponse            // Update response for the next iteration
		}

		return nil
	})

	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return s.ctx.Err()
	default:
		if err := g.Wait(); err != nil {
			return err
		}

		return nil
	}
}

func (s *Stream) Unsubscribe() error {
	close(s.ch)
	close(s.errCh)
	s.cancelFn()
	return nil
}

func (s *Stream) QueueError(err error) {
	s.errCh <- err
}

func (s *Stream) Err() <-chan error {
	return s.errCh
}

func (s *Stream) Channel() <-chan *types.QueryResponse {
	return s.ch
}

func (s *Stream) Done() <-chan struct{} {
	return s.done
}

func (s *Stream) ChannelWithError() (<-chan *types.QueryResponse, <-chan error) {
	return s.ch, s.errCh
}
