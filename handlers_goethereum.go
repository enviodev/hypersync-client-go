package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/types"
	"math/big"
)

func (c *Client) HeaderByNumber(ctx context.Context, blockNumber *big.Int) (*types.QueryResponse, error) {
	query := types.Query{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		FieldSelection: types.FieldSelection{
			Block: []string{"number", "hash"},
		},
	}
	return c.GetArrow(ctx, &query)
}
