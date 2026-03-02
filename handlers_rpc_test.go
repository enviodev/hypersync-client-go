package hypersyncgo

import (
	"context"
	"math/big"
	"testing"

	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHeaderByNumber(t *testing.T) {
	skipWithoutApiToken(t)
	testCases := []struct {
		name         string
		opts         options.Options
		blockNumbers []*big.Int
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
					ApiToken:    getTestApiToken(),
				},
			},
		},
		blockNumbers: []*big.Int{
			big.NewInt(10000000),
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

			for _, q := range testCase.blockNumbers {
				resp, rErr := client.HeaderByNumber(ctx, q)
				require.NoError(t, rErr)
				require.NotNil(t, resp)
				require.Equal(t, resp.Number, q)
			}
		})
	}
}

func TestBlockByNumber(t *testing.T) {
	skipWithoutApiToken(t)
	testCases := []struct {
		name         string
		opts         options.Options
		blockNumbers []*big.Int
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
					ApiToken:    getTestApiToken(),
				},
			},
		},
		blockNumbers: []*big.Int{
			big.NewInt(10000000),
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

			for _, q := range testCase.blockNumbers {
				resp, rErr := client.BlockByNumber(ctx, q)
				require.NoError(t, rErr)
				require.NotNil(t, resp)
			}
		})
	}
}
