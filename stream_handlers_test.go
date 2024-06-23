package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/logger"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"math/big"
	"testing"
	"time"
)

func TestGetBlocksInRange(t *testing.T) {
	// Here just to test somewhere that logger is actually loaded...
	zLog, err := logger.GetZapDevelopmentLogger(zap.NewAtomicLevelAt(zap.DebugLevel))
	require.NoError(t, err)
	require.NotNil(t, zLog)
	logger.SetGlobalLogger(zLog)

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
				bStream, bsErr := client.StreamInRange(ctx, r.start, r.end, r.options)
				require.Nil(t, bsErr)
				require.NotNil(t, bStream)

				select {
				case cErr := <-bStream.Err():
					t.Errorf("Got error from GetBlocksInRange: %s", cErr)
					require.Nil(t, bStream.Unsubscribe())
					require.NotNil(t, cErr)
				case response := <-bStream.Channel():
					utils.DumpNodeNoExit(response)
				case <-time.After(2 * time.Second):
					require.Nil(t, bStream.Unsubscribe())
					require.Fail(t, "expected ranges to contain at least one block")
				}
			}
		})
	}
}
