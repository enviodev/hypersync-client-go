//go:build ignore
// +build ignore

package main

import (
	"context"
	"math/big"
	"time"

	hypersyncgo "github.com/enviodev/hypersync-client-go"
	"github.com/enviodev/hypersync-client-go/logger"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

func main() {
	opts := options.Options{
		Blockchains: []options.Node{
			{
				Type:        utils.EthereumNetwork,
				NetworkId:   utils.EthereumNetworkID,
				Endpoint:    "https://arbitrum.hypersync.xyz",
				RpcEndpoint: "https://arbitrum.rpc.hypersync.xyz",
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hsClient, err := hypersyncgo.NewHyper(ctx, opts)
	if err != nil {
		logger.L().Error(
			"failed to create hyper client",
			zap.Error(err),
		)
		return
	}

	client, found := hsClient.GetClient(utils.EthereumNetworkID)
	if !found {
		logger.L().Error(
			"failure to discover hyper client",
			zap.Error(err),
			zap.Any("network_id", utils.EthereumNetworkID),
		)
		return
	}

	startBlock := big.NewInt(274569854)
	endBlock := big.NewInt(274569857)
	startTime := time.Now()

	logger.L().Info(
		"New signature hash with custom stream query request started",
		zap.Error(err),
		zap.Any("network_id", utils.EthereumNetworkID),
		zap.Any("start_block", startBlock),
		zap.Any("end_block", endBlock),
	)

	query := types.Query{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Transactions: []types.TransactionSelection{
			{
				SigHash: []types.SigHash{
					types.NewSigHashFromHex("ad718d2a"),
				},
			},
		},
		Logs: []types.LogSelection{
			{
				Topics: [][]common.Hash{
					{common.HexToHash("0x7c3072652d5407e4dcbe90ac1760509311e2511e531f5caf91859f3bc7416708")},
				},
			},
		},
		FieldSelection: types.FieldSelection{
			Block:       []string{"number", "timestamp", "hash"},
			Transaction: []string{"from", "to", "value", "input", "block_number", "block_hash", "hash", "contract_address", "sighash"},
			Log:         []string{"address", "data", "topic0", "topic1", "topic2", "topic3", "block_number", "block_hash", "transaction_hash"},
		},
	}

	batchSize := big.NewInt(50)
	bStream, bsErr := client.Stream(ctx, &query, options.DefaultStreamOptionsWithBatchSize(batchSize))
	if bsErr != nil {
		logger.L().Error(
			"failure to execute hyper client stream in range",
			zap.Error(err),
			zap.Any("network_id", utils.EthereumNetworkID),
			zap.Any("start_block", startBlock),
			zap.Any("end_block", endBlock),
		)
		return
	}

	latestBatchReceived := big.NewInt(0)
	totalBlocks := make(map[uint64]struct{})
	totalTxns := 0
	totalLogs := 0
	for {
		select {
		case cErr := <-bStream.Err():
			logger.L().Error(
				"failure to execute hyper client stream in range",
				zap.Error(cErr),
				zap.Any("network_id", utils.EthereumNetworkID),
				zap.Any("start_block", startBlock),
				zap.Any("end_block", endBlock),
			)
			return
		case response := <-bStream.Channel():
			logger.L().Info(
				"New stream logs response",
				zap.Any("start_block", startBlock),
				zap.Any("current_sync_block", response.NextBlock),
				zap.Any("end_block", endBlock),
				zap.Duration("current_processing_time", time.Since(startTime)),
			)
			latestBatchReceived = response.NextBlock

			totalTxns += len(response.GetTransactions())
			totalLogs += len(response.GetLogs())

			for _, block := range response.GetBlocks() {
				totalBlocks[block.Number.Uint64()] = struct{}{}
				logger.L().Info(
					"Block information",
					zap.Any("number", block.Number),
					zap.Any("hash", block.Hash),
					zap.Any("timestamp", block.Timestamp),
				)
			}

			// Worker may close a done channel at any point in time after it receives all the payload.
			// Usually and depending on the payload size,
			// it takes around 5-100ms after stream is completed for messages to fully be delivered.
			// Instead of using time.Sleep(), we have Ack() mechanism that allows you finer worker closure management.
			// WARN: This is critical part of communication with stream and should always be used unless you have
			// disabled it via configuration. By default, ack is a must!
			bStream.Ack()

		case <-bStream.Done():
			logger.L().Info(
				"Stream request successfully completed",
				zap.Duration("duration", time.Since(startTime)),
				zap.Any("total_blocks", len(totalBlocks)),
				zap.Any("total_txns", totalTxns),
				zap.Any("total_logs", totalLogs),
			)
			return
		case <-time.After(15 * time.Second):
			logger.L().Error(
				"expected ranges to receive at least one logs range in 15s",
				zap.Any("network_id", utils.EthereumNetworkID),
				zap.Any("start_block", startBlock),
				zap.Any("latest_batch_block_received", latestBatchReceived),
				zap.Any("end_block", endBlock),
			)
			return
		}
	}
}
