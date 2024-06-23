package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestClients(t *testing.T) {
	//toBlock := uint64(10000001)

	testCases := []struct {
		name    string
		opts    options.Options
		queries []*types.Query
	}{{
		name: "Test Ethereum Client",
		opts: options.Options{
			Blockchains: []options.Node{
				{
					Type:      utils.EthereumNetwork,
					NetworkId: utils.EthereumNetworkID,
					Endpoint:  "https://eth.hypersync.xyz",
				},
			},
		},
		queries: []*types.Query{
			{
				FromBlock: big.NewInt(10000000),
				Transactions: []types.TransactionSelection{
					{
						ContractAddress: []types.Address{
							"0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE",
						},
					},
				},
				FieldSelection: types.FieldSelection{
					Block:       []string{"number", "hash"},
					Transaction: []string{"hash", "block_number"},
				},
			},
		},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Fetch the first node out of the blockchain definitions
			nodeOpts := testCase.opts.Blockchains[0]
			client, err := NewClient(ctx, nodeOpts)
			require.NoError(t, err)
			require.NotNil(t, client)

			height, err := client.GetHeight(ctx)
			require.NoError(t, err)
			t.Logf("Discovered current height: %d", height)
			require.Greater(t, height, uint64(0))

			for _, q := range testCase.queries {
				resp, rErr := client.Get(ctx, q)
				require.NoError(t, rErr)
				require.NotNil(t, resp)
				utils.DumpNodeNoExit(resp)
			}

		})
	}
}
