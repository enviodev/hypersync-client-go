package client

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go/types"
	"math/rand"
	"net/http"
	"time"
)

func (c *Client) GetArrow(ctx context.Context, query *types.Query) (*types.QueryResponse[interface{}], error) {
	base := c.opts.RetryBaseMs

	c.opts.RetryBackoffMs = time.Duration(100)
	c.opts.MaxNumRetries = 0

	for i := 0; i < c.opts.MaxNumRetries+1; i++ {
		response, err := DoArrow[*types.Query, types.QueryResponse[interface{}]](ctx, c, c.GeUrlFromNodeAndPath(c.opts, "query", "arrow-ipc"), http.MethodPost, query)
		if err == nil {
			return response, nil
		}

		// TODO: Implement proper logger...
		fmt.Printf("Failed to get arrow data from server, retrying... Error: %v\n", err)

		baseMs := base * time.Millisecond
		jitter := time.Duration(rand.Int63n(int64(c.opts.RetryBackoffMs))) * time.Millisecond

		select {
		case <-time.After(baseMs + jitter):
			base = min(base+c.opts.RetryBackoffMs, c.opts.RetryCeilingMs)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("failed to get arrow data after retries: %d", c.opts.MaxNumRetries)
}
