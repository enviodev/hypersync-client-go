package main

/*
 This is a fun example of decoding ERC721 events with the built-in decoder.  We use the Clovers.network contracts as an example, which are tokens that are minted by players who work out a symetric othelo game. View more at https://clovers.network
*/
import (
	"context"
	"fmt"
	"math/big"
	"time"

	hypersyncgo "github.com/enviodev/hypersync-client-go"
	"github.com/enviodev/hypersync-client-go/contracts"
	"github.com/enviodev/hypersync-client-go/decoder"
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
				Endpoint:    "https://eth.hypersync.xyz",
				RpcEndpoint: "https://eth.rpc.hypersync.xyz",
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hsClient, err := hypersyncgo.NewHyper(ctx, opts)
	if err != nil {
		logger.L().Error("failed to create hyper client", zap.Error(err))
		return
	}

	client, found := hsClient.GetClient(utils.EthereumNetworkID)
	if !found {
		logger.L().Error("failure to discover hyper client", zap.Error(err))
		return
	}

	startBlock := big.NewInt(20000000)
	endBlock := big.NewInt(20865800)

	logger.L().Info("Fetching USDC transfer logs", zap.Any("start_block", startBlock), zap.Any("end_block", endBlock))

	// Clovers Transfer event signature: Transfer(address,address,uint256)
	selections := []types.LogSelection{
		{
			Address: []common.Address{common.HexToAddress("0xB55C5cAc5014C662fDBF21A2C59Cd45403C482Fd")}, // Clovers contract address
			Topics: [][]common.Hash{
				{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}, // Transfer event signature
			},
		},
	}

	batchSize := big.NewInt(50)
	bStream, bsErr := client.StreamLogsInRange(ctx, startBlock, endBlock, selections, options.DefaultStreamOptionsWithBatchSize(batchSize))
	if bsErr != nil {
		logger.L().Error("failure to execute hyper client stream in range", zap.Error(bsErr))
		return
	}

	// Load contracts manager to decode logs
	cManager, cErr := contracts.NewManagerWithDefaults()
	if cErr != nil {
		logger.L().Error("failed to load contracts manager", zap.Error(cErr))
		return
	}

	contract, cdErr := cManager.GetByStandard(utils.EthereumNetworkID, utils.Erc721)
	if cdErr != nil {
		logger.L().Error("failed to get ERC721 interface for Clovers contract", zap.Error(cdErr))
		return
	}

	for {
		select {
		case cErr := <-bStream.Err():
			logger.L().Error("stream error", zap.Error(cErr))
			return
		case response := <-bStream.Channel():
			for _, log := range response.GetLogs() {
				dLog, dlErr := decoder.DecodeEthereumLogWithContract(log, contract)
				if dlErr != nil {
					logger.L().Error("failed to decode log", zap.Error(dlErr))
					continue
				}

				fmt.Printf("Decoded Clovers event: %+v\n", dLog)
			}
			bStream.Ack()
		case <-bStream.Done():
			logger.L().Info("Stream request successfully completed")
			return
		case <-time.After(15 * time.Second):
			logger.L().Error("Timeout: no logs received in the last 15s")
			return
		}
	}
}
