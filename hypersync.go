package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/pkg/options"
	"github.com/pkg/errors"
)

type HyperSync struct {
	ctx  context.Context
	opts options.Options
}

func NewHyperSync(ctx context.Context, opts options.Options) (*HyperSync, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid options to hypersync client")
	}
	return &HyperSync{
		ctx:  ctx,
		opts: opts,
	}, nil
}
