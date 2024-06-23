package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestGetContractCreatorByNumber(t *testing.T) {
	testCases := []struct {
		name  string
		opts  options.Options
		cases []struct {
			address common.Address
			number  *big.Int
			hash    common.Hash
		}
	}{{
		name: "Test Ethereum Client",
		opts: options.Options{
			Blockchains: []options.Node{
				{
					Type:        utils.EthereumNetwork,
					NetworkId:   utils.EthereumNetworkID,
					Endpoint:    "https://eth.hypersync.xyz",
					RpcEndpoint: "https://eth.rpc.hypersync.xyz",
				},
			},
		},
		cases: []struct {
			address common.Address
			number  *big.Int
			hash    common.Hash
		}{
			{
				address: common.HexToAddress("0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE"),
				number:  big.NewInt(10569013),
				hash:    common.HexToHash("0x678dc99c448fc2dbc10081160066b5f654c916340c79c3b239fe4aaad200dca9"),
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

			for _, q := range testCase.cases {
				resp, rErr := client.GetContractCreator(ctx, q.address)
				require.NoError(t, rErr)
				require.NotNil(t, resp)
				require.Equal(t, resp.Number, q.number)
				require.Equal(t, resp.Hash, q.hash)
			}
		})
	}
}
