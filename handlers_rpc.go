package hypersyncgo

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func (c *Client) HeaderByNumber(ctx context.Context, blockNumber *big.Int) (*types.Header, error) {
	return c.GetRPC().HeaderByNumber(ctx, blockNumber)
}

func (c *Client) BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error) {
	return c.GetRPC().BlockByNumber(ctx, blockNumber)
}
