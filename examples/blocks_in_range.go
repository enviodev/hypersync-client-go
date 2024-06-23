package main

import (
	"context"
	"fmt"
	"github.com/enviodev/hypersync-client-go"
	"github.com/enviodev/hypersync-client-go/logger"
	"github.com/enviodev/hypersync-client-go/options"
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

	startBlock := big.NewInt(10000000)
	endBlock := big.NewInt(10015000)

	logger.L().Info(
		"New blocks in range stream request started",
		zap.Error(err),
		zap.Any("network_id", utils.EthereumNetworkID),
		zap.Any("start_block", startBlock),
		zap.Any("end_block", endBlock),
	)

	bStream, bsErr := client.StreamInRange(ctx, startBlock, endBlock, options.DefaultStreamOptions())
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
				zap.Any("next_block", response.NextBlock),
			)

			for _, block := range response.Data.Blocks {
				fmt.Println("Hash:", block.Hash)
			}

		case <-bStream.Done():
			logger.L().Info("Stream request successfully completed")
			return
		case <-time.After(5 * time.Second):
			logger.L().Error(
				"expected ranges to receive at least one block in 5s",
				zap.Any("network_id", utils.EthereumNetworkID),
				zap.Any("start_block", startBlock),
				zap.Any("end_block", endBlock),
			)
			return
		}
	}
}