package hypersyncgo

import (
	"github.com/ethereum/go-ethereum/common"
)

// ContractCreatorFieldSelection represents the fields to be selected from the transaction response.
type ContractCreatorFieldSelection struct {
	Transaction []string `json:"transaction"`
}

// ContractCreatorTransaction represents the transaction data containing contract addresses.
type ContractCreatorTransaction struct {
	ContractAddress []common.Address `json:"contract_address"`
}

// ContractCreatorRequest represents the request structure for fetching contract creator details.
type ContractCreatorRequest struct {
	FromBlock      int                           `json:"from_block"`
	Transactions   []ContractCreatorTransaction  `json:"transactions"`
	FieldSelection ContractCreatorFieldSelection `json:"field_selection"`
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
/*func (c *Client) GetContractCreator(ctx context.Context, addr common.Address) (*Transaction, error) {
	request := ContractCreatorRequest{
		FromBlock: 0,
		Transactions: []ContractCreatorTransaction{
			{
				ContractAddress: []common.Address{addr},
			},
		},
		FieldSelection: ContractCreatorFieldSelection{
			Transaction: []string{"hash", "block_number"},
		},
	}

	response, err := DoQuery[ContractCreatorRequest, Response](ctx, c, http.MethodPost, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get contract creator from envio")
	}

	if len(response.Data) == 0 {
		return nil, errorshs.ErrContractNotFound
	}

	return response.Data[0].Transactions[0], nil
}*/
