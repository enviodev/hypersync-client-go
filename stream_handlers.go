package hypersyncgo

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"math/big"
)

func (c *Client) StreamBlocksInRange(ctx context.Context, fromBlock *big.Int, toBlock *big.Int, opts *options.StreamOptions) (*Stream, error) {
	if fromBlock == nil {
		return nil, fmt.Errorf("from block must not be nil")
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
			Block: types.BlockSchemaFieldsAsString(),
		},
	}

	if opts == nil {
		opts = options.DefaultStreamOptions()
	}

	return c.Stream(ctx, &query, opts)
}

func (c *Client) StreamTransactionsInRange(ctx context.Context, fromBlock *big.Int, toBlock *big.Int, selections []types.TransactionSelection, opts *options.StreamOptions) (*Stream, error) {
	if fromBlock == nil {
		return nil, fmt.Errorf("from block must not be nil")
	}

	// Querying will not return toBlock actual value that was requested but rather toBlock-1
	if toBlock != nil && toBlock.Cmp(big.NewInt(0)) > 0 {
		toBlock = toBlock.Add(toBlock, big.NewInt(1))
	}

	query := types.Query{
		FromBlock:        fromBlock,
		IncludeAllBlocks: true, // We have to include all blocks otherwise data won't be returned back.
		ToBlock:          toBlock,
		Transactions:     selections,
		FieldSelection: types.FieldSelection{
			Transaction: types.TransactionSchemaFieldsAsString(),
		},
	}

	if opts == nil {
		opts = options.DefaultStreamOptions()
	}

	return c.Stream(ctx, &query, opts)
}
