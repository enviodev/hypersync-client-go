package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"math/big"
	"testing"
)

func TestClients(t *testing.T) {
	testCases := []struct {
		name    string
		opts    options.Options
		queries []*types.Query
	}{{
		name: "Test Ethereum Client",
		opts: options.Options{
			LogLevel: zap.DebugLevel,
			Blockchains: []options.Node{
				{
					Type:        utils.EthereumNetwork,
					NetworkId:   utils.EthereumNetworkID,
					Endpoint:    "https://eth.hypersync.xyz",
					RpcEndpoint: "https://eth.rpc.hypersync.xyz",
				},
			},
		},
		queries: []*types.Query{
			{
				FromBlock: big.NewInt(10000000),
				Transactions: []types.TransactionSelection{
					{
						ContractAddress: []common.Address{
							common.HexToAddress("0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE"),
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

			for _, q := range testCase.queries {
				resp, rErr := client.Get(ctx, q)
				require.NoError(t, rErr)
				require.NotNil(t, resp)
			}
		})
	}
}
