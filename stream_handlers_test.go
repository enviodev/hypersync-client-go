package hypersyncgo

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/markovichecha/hypersync-client-go/logger"
	"github.com/markovichecha/hypersync-client-go/options"
	"github.com/markovichecha/hypersync-client-go/types"
	"github.com/markovichecha/hypersync-client-go/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
				end:     big.NewInt(10020000),
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
				bStream, bsErr := client.StreamBlocksInRange(ctx, r.start, r.end, r.options)
				require.Nil(t, bsErr)
				require.NotNil(t, bStream)

				for {
					select {
					case cErr := <-bStream.Err():
						t.Errorf("Got error from GetBlocksInRange: %s", cErr)
						//require.Nil(t, bStream.Unsubscribe())
						require.NotNil(t, cErr)
					case <-bStream.Done():
						t.Log("Stream successfully resolved!")
						return
					case response := <-bStream.Channel():
						t.Logf("Got response from GetBlocksInRange NextBlock: %d", response.NextBlock)
						//utils.DumpNodeNoExit(response)
						bStream.Ack()
					case <-time.After(5 * time.Second):
						//require.Nil(t, bStream.Unsubscribe())
						require.Fail(t, "expected ranges to receive at least one block in 5s")
					}
				}
			}
		})
	}
}

func TestGetTransactionsInRange(t *testing.T) {
	// Here just to test somewhere that logger is actually loaded...
	zLog, err := logger.GetZapDevelopmentLogger(zap.NewAtomicLevelAt(zap.DebugLevel))
	require.NoError(t, err)
	require.NotNil(t, zLog)
	logger.SetGlobalLogger(zLog)

	statusValue := uint8(1)

	testCases := []struct {
		name      string
		opts      options.Options
		networkId utils.NetworkID
		ranges    []struct {
			start      *big.Int
			end        *big.Int
			selections []types.TransactionSelection
			options    *options.StreamOptions
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
			start      *big.Int
			end        *big.Int
			selections []types.TransactionSelection
			options    *options.StreamOptions
		}{
			{
				start: big.NewInt(10000000),
				end:   big.NewInt(10000200),
				selections: []types.TransactionSelection{
					{
						Status: &statusValue,
					},
				},
				options: options.DefaultStreamOptionsWithBatchSize(big.NewInt(10)),
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
				bStream, bsErr := client.StreamTransactionsInRange(ctx, r.start, r.end, r.selections, r.options)
				require.Nil(t, bsErr)
				require.NotNil(t, bStream)

				for {
					select {
					case cErr := <-bStream.Err():
						t.Errorf("Got error from StreamTransactionsInRange: %s", cErr)
						require.NotNil(t, cErr)
					case <-bStream.Done():
						t.Log("Stream successfully resolved!")
						//require.Nil(t, bStream.Unsubscribe())
						return
					case response := <-bStream.Channel():
						t.Logf("Got response from StreamTransactionsInRange NextBlock: %d", response.NextBlock)
						bStream.Ack()
					case <-time.After(15 * time.Second):
						//require.Nil(t, bStream.Unsubscribe())
						require.Fail(t, "expected ranges to receive at least one block in 15s")
					}
				}
			}
		})
	}
}

func TestGetLogsInRange(t *testing.T) {
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
			start      *big.Int
			end        *big.Int
			selections []types.LogSelection
			options    *options.StreamOptions
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
			start      *big.Int
			end        *big.Int
			selections []types.LogSelection
			options    *options.StreamOptions
		}{
			{
				start: big.NewInt(20000000),
				end:   big.NewInt(20000100),
				selections: []types.LogSelection{
					{
						Topics: [][]common.Hash{
							{
								// Transfer(address,address,uint256)
								common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
							},
						},
					},
				},
				options: options.DefaultStreamOptionsWithBatchSize(big.NewInt(10)),
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
				bStream, bsErr := client.StreamLogsInRange(ctx, r.start, r.end, r.selections, r.options)
				require.Nil(t, bsErr)
				require.NotNil(t, bStream)

				for {
					select {
					case cErr := <-bStream.Err():
						t.Errorf("Got error from StreamLogsInRange: %s", cErr)
						require.NotNil(t, cErr)
					case <-bStream.Done():
						t.Log("Stream successfully resolved!")
						//require.Nil(t, bStream.Unsubscribe())
						return
					case response := <-bStream.Channel():
						t.Logf("Got response from StreamLogsInRange NextBlock: %d", response.NextBlock)
						bStream.Ack()
					case <-time.After(15 * time.Second):
						//require.Nil(t, bStream.Unsubscribe())
						require.Fail(t, "expected ranges to receive at least one block in 15s")
					}
				}
			}
		})
	}
}

func TestGetTracesInRange(t *testing.T) {
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
			start      *big.Int
			end        *big.Int
			selections []types.TraceSelection
			options    *options.StreamOptions
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
			start      *big.Int
			end        *big.Int
			selections []types.TraceSelection
			options    *options.StreamOptions
		}{
			{
				start: big.NewInt(20000000),
				end:   big.NewInt(20000100),
				selections: []types.TraceSelection{
					{
						To: []common.Address{
							// Uniswap V2: Router 2
							common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"),
						},
					},
				},
				options: options.DefaultStreamOptionsWithBatchSize(big.NewInt(10)),
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
				bStream, bsErr := client.StreamTracesInRange(ctx, r.start, r.end, r.selections, r.options)
				require.Nil(t, bsErr)
				require.NotNil(t, bStream)

				for {
					select {
					case cErr := <-bStream.Err():
						t.Errorf("Got error from StreamTracesInRange: %s", cErr)
						require.NotNil(t, cErr)
					case <-bStream.Done():
						t.Log("Stream successfully resolved!")
						//require.Nil(t, bStream.Unsubscribe())
						return
					case response := <-bStream.Channel():
						t.Logf("Got response from StreamTracesInRange NextBlock: %d", response.NextBlock)
						bStream.Ack()
					case <-time.After(15 * time.Second):
						//require.Nil(t, bStream.Unsubscribe())
						require.Fail(t, "expected ranges to receive at least one block in 15s")
					}
				}
			}
		})
	}
}
