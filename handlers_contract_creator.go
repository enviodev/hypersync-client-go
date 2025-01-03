package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

type ContractCreatorResponse struct {
	Number *big.Int    `json:"number"`
	Hash   common.Hash `json:"hash"`
	From common.Address `json:"from"`
}

// GetContractCreator fetches the transaction details of the creator of a specified contract address.
//
// This function sends a POST request to the Envio API to retrieve information about the contract creator.
// It takes the context, network ID, and contract address as parameters and returns the transaction details.
//
// Parameters:
// - ctx: The context for managing the request lifecycle.
// - networkId: The ID of the network where the contract is deployed.
// - addr: The contract address for which the creator's transaction details are to be fetched.
//
// Returns:
// - *Transaction: The transaction details of the contract creator, if found.
// - error: An error if the request fails or the contract is not found.
func (c *Client) GetContractCreator(ctx context.Context, addr common.Address) (*ContractCreatorResponse, error) {
	query := types.Query{
		FromBlock: big.NewInt(0),
		Transactions: []types.TransactionSelection{
			{
				ContractAddress: []common.Address{addr},
			},
		},
		FieldSelection: types.FieldSelection{
			Block:       []string{"number"},
			Transaction: []string{"hash", "from"},
		},
	}

	response, err := c.GetArrow(ctx, &query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get contract creator by address: %s", addr)
	}

	return &ContractCreatorResponse{
		Number: response.Data.Blocks[0].Number,
		Hash:   *response.Data.Transactions[0].Hash,
		From: 	*response.Data.Transactions[0].From,
	}, nil
}
