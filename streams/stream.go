package streams

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Stream struct {
	ctx      context.Context
	cancelFn context.CancelFunc
	ch       chan *types.QueryResponse
	errCh    chan error
	opts     *options.StreamOptions
	query    *types.Query
}

func NewStream(ctx context.Context, query *types.Query, opts *options.StreamOptions) (*Stream, error) {
	ctx, cancel := context.WithCancel(ctx)
	return &Stream{
		ctx:      ctx,
		cancelFn: cancel,
		query:    query,
		ch:       make(chan *types.QueryResponse, opts.Concurrency.Uint64()),
		errCh:    make(chan error, opts.Concurrency.Uint64()),
	}, nil
}

func (s *Stream) Subscribe() error {
	g, ctx := errgroup.WithContext(s.ctx)

	g.Go(func() error {
		return errors.New("DUMMY_TEST_ERROR")
	})

	select {
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

func (s *Stream) ChannelWithError() (<-chan *types.QueryResponse, <-chan error) {
	return s.ch, s.errCh
}
