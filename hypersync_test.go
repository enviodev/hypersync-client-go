package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHyperSync(t *testing.T) {
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
					Type:      utils.EthereumNetwork,
					NetworkId: utils.EthereumNetworkID,
					Endpoint:  "https://eth.hypersync.xyz",
				},
			},
		},
		networkId: utils.EthereumNetworkID,
		addresses: []common.Address{
			common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
		},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			hsClient, err := NewHyperSync(ctx, testCase.opts)
			require.NoError(t, err)
			require.NotNil(t, hsClient)

			client, found := hsClient.GetClient(testCase.networkId)
			require.NoError(t, err)
			require.True(t, found)

			height, err := client.GetHeight(ctx)
			require.NoError(t, err)
			t.Logf("Discovered current height: %d", height)
			require.Greater(t, height, uint64(0))
		})
	}
}
