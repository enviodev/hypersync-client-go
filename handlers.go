package hypersyncgo

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go/types"
	"math/big"
	"math/rand"
	"net/http"
	"time"
)

type ArchiveHeight struct {
	Height *big.Int `json:"height"`
}

func (c *Client) GetHeight(ctx context.Context) (*big.Int, error) {
	base := c.opts.RetryBaseMs

	for i := 0; i < c.opts.MaxNumRetries+1; i++ {
		response, err := Do[struct{}, ArchiveHeight](ctx, c, c.GeUrlFromNodeAndPath(c.opts, "height"), http.MethodGet, struct{}{})
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
			return big.NewInt(0), ctx.Err()
		}
	}

	return big.NewInt(0), fmt.Errorf("failed to get height after retries: %d", c.opts.MaxNumRetries)
}

func (c *Client) Get(ctx context.Context, query *types.Query) (*types.QueryResponse, error) {
	return c.GetArrow(ctx, query)
}

func (c *Client) GetBlocksInRange(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) (*types.QueryResponse, error) {
	if fromBlock == nil {
		return nil, fmt.Errorf("fromBlock must not be nil")
	}

	// Querying will not return toBlock actual value that was requested but rather toBlock-1
	if toBlock != nil && toBlock.Cmp(big.NewInt(0)) > 0 {
		toBlock = toBlock.Add(toBlock, big.NewInt(1))
	}

	query := types.Query{
		FromBlock:        fromBlock,
		IncludeAllBlocks: true, // We have to include all blocks otherwise data won't be returned back.
		ToBlock:          toBlock,
		FieldSelection: types.FieldSelection{
			Block: []string{"number", "hash"},
		},
	}

	return c.GetArrow(ctx, &query)
}
