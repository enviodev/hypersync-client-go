//go:build ignore
// +build ignore

package main

import (
	"context"
	"github.com/enviodev/hypersync-client-go"
	"github.com/enviodev/hypersync-client-go/logger"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/enviodev/hypersync-client-go/utils"
	"go.uber.org/zap"
	"math/big"
	"time"
)

func main() {
	opts := options.Options{
		Blockchains: []options.Node{
			{
				Type:        utils.EthereumNetwork,
				NetworkId:   utils.EthereumNetworkID,
				Endpoint:    "https://eth.hypersync.xyz",
				RpcEndpoint: "https://eth.rpc.hypersync.xyz",
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

	startBlock := big.NewInt(20000000)
	endBlock := big.NewInt(20001000)
	startTime := time.Now()

	logger.L().Info(
		"New blocks in range stream request started",
		zap.Error(err),
		zap.Any("network_id", utils.EthereumNetworkID),
		zap.Any("start_block", startBlock),
		zap.Any("end_block", endBlock),
	)

	statusValue := uint8(1)
	selections := []types.TransactionSelection{
		{
			Status: &statusValue,
		},
	}

	batchSize := big.NewInt(50)
	bStream, bsErr := client.StreamTransactionsInRange(ctx, startBlock, endBlock, selections, options.DefaultStreamOptionsWithBatchSize(batchSize))
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
				"New stream block response",
				zap.Any("start_block", startBlock),
				zap.Any("current_sync_block", response.NextBlock),
				zap.Any("end_block", endBlock),
				zap.Duration("current_processing_time", time.Since(startTime)),
			)
			latestBatchReceived = response.NextBlock

			totalTxns += len(response.GetTransactions())
			for _, tx := range response.GetTransactions() {
				totalBlocks[tx.BlockNumber.Uint64()] = struct{}{}
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
				zap.Any("total_transactions", totalTxns),
			)
			return
		case <-time.After(15 * time.Second):
			logger.L().Error(
				"expected ranges to receive at least one block range in 15s",
				zap.Any("network_id", utils.EthereumNetworkID),
				zap.Any("start_block", startBlock),
				zap.Any("latest_batch_block_received", latestBatchReceived),
				zap.Any("end_block", endBlock),
			)
			return
		}
	}
}
