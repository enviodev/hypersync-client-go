package hypersyncgo

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go/types"
	"math/rand"
	"net/http"
	"time"
)

func (c *Client) GetHeight(ctx context.Context) (uint64, error) {
	base := c.opts.RetryBaseMs

	for i := 0; i < c.opts.MaxNumRetries+1; i++ {
		response, err := Do[struct{}, types.ArchiveHeight](ctx, c, c.GeUrlFromNodeAndPath(c.opts, "height"), http.MethodGet, struct{}{})
		if err == nil {
			return response.Height, nil
		}

		// TODO: Implement proper logger...
		fmt.Printf("Failed to get height from server, retrying... Error: %v\n", err)

		baseMs := base * time.Millisecond
		jitter := time.Duration(rand.Int63n(int64(c.opts.RetryBackoffMs))) * time.Millisecond

		select {
		case <-time.After(baseMs + jitter):
			base = min(base+c.opts.RetryBackoffMs, c.opts.RetryCeilingMs)
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	return 0, fmt.Errorf("failed to get height after retries: %d", c.opts.MaxNumRetries)
}

func (c *Client) Get(ctx context.Context, query *types.Query) (*types.QueryResponse, error) {
	return c.GetArrow(ctx, query)
}
