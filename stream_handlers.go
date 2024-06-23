package hypersyncgo

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"math/big"
)

func (c *Client) GetBlocksInRange(ctx context.Context, fromBlock *big.Int, toBlock *big.Int, opts *options.StreamOptions) (<-chan *types.QueryResponse, error) {
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

	if opts == nil {
		opts = options.DefaultStreamOptions()
	}

	return c.Stream(ctx, &query, opts)
}
