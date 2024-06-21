package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/pkg/options"
	"github.com/enviodev/hypersync-client-go/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClients(t *testing.T) {
	testCases := []struct {
		name      string
		opts      options.Options
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
		addresses: []common.Address{
			common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
		},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client, err := NewClient(ctx, testCase.opts)
			require.NoError(t, err)
			require.NotNil(t, client)
		})
	}
}
