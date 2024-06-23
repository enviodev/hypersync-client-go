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

func TestGetHeight(t *testing.T) {
	testCases := []struct {
		name      string
		opts      options.Options
		networkId utils.NetworkID
		addresses []common.Address
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

			height, err := client.GetHeight(ctx)
			require.NoError(t, err)
			t.Logf("Discovered current height: %d", height)
			require.Greater(t, height.Uint64(), big.NewInt(0).Uint64())
		})
	}
}

func TestGetBlocksInRange(t *testing.T) {
	testCases := []struct {
		name      string
		opts      options.Options
		networkId utils.NetworkID
		ranges    []struct {
			start *big.Int
			end   *big.Int
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
			start *big.Int
			end   *big.Int
		}{
			{
				start: big.NewInt(10000000),
				end:   big.NewInt(10001002),
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
				ranges, rErr := client.GetBlocksInRange(ctx, r.start, r.end)
				require.NoError(t, rErr)
				require.NotNil(t, ranges)
				utils.DumpNodeNoExit(ranges)
			}
		})
	}
}
