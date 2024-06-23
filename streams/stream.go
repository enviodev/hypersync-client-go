package streams

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
)

type Stream struct {
	ctx   context.Context
	ch    chan *types.QueryResponse
	opts  *options.StreamOptions
	query *types.Query
}

func NewStream(ctx context.Context, query *types.Query, opts *options.StreamOptions) (*Stream, error) {
	return &Stream{
		ctx:   ctx,
		query: query,
		ch:    make(chan *types.QueryResponse, opts.Concurrency.Uint64()),
	}, nil
}

func (s *Stream) Channel() <-chan *types.QueryResponse {
	return s.ch
}
