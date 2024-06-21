package hypersyncgo

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go/pkg/types"
	"math/rand"
	"net/http"
	"time"
)

func (c *Client) GetHeight(ctx context.Context) (uint64, error) {
	base := c.nodeOpts.RetryBaseMs

	for i := 0; i < c.nodeOpts.MaxNumRetries+1; i++ {
		response, err := Do[struct{}, types.ArchiveHeight](ctx, c, c.GeUrlFromNodeAndPath(c.nodeOpts, "height"), http.MethodGet, struct{}{})
		if err == nil {
			return response.Height, nil
		}

		fmt.Printf("Failed to get height from server, retrying... Error: %v\n", err)

		baseMs := base * time.Millisecond
		jitter := time.Duration(rand.Int63n(int64(c.nodeOpts.RetryBackoffMs))) * time.Millisecond

		select {
		case <-time.After(baseMs + jitter):
			base = min(base+c.nodeOpts.RetryBackoffMs, c.nodeOpts.RetryCeilingMs)
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	return 0, fmt.Errorf("failed to get height after retries")
}
