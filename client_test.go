package hypersyncgo

import (
	"context"
	"github.com/enviodev/hypersync-client-go/pkg/options"
	"github.com/enviodev/hypersync-client-go/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClients(t *testing.T) {
	testCases := []struct {
		name string
		opts options.Options
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
