package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestGetBlocksInRange(t *testing.T) {
	testCases := []struct {
		name      string
		opts      options.Options
		networkId utils.NetworkID
		ranges    []struct {
			start   *big.Int
			end     *big.Int
			options *options.StreamOptions
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
		networkId: utils.EthereumNetworkID,
		ranges: []struct {
			start   *big.Int
			end     *big.Int
			options *options.StreamOptions
		}{
			{
				start:   big.NewInt(10000000),
				end:     big.NewInt(11000000),
				options: options.DefaultStreamOptions(),
			},
		},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			hsClient, err := NewHyper(ctx, testCase.opts)
			require.NoError(t, err)
			require.NotNil(t, hsClient)

			client, found := hsClient.GetClient(testCase.networkId)
			require.True(t, found)
			require.NotNil(t, client)

			for _, r := range testCase.ranges {
				ranges, rErr := client.GetBlocksInRange(ctx, r.start, r.end, r.options)
				require.NoError(t, rErr)
				require.NotNil(t, ranges)

				select {
				case <-ranges:

				case <-time.After(2 * time.Second):
					require.Fail(t, "expected ranges to contain at least one block")
				}
			}
		})
	}
}