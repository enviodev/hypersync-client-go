package hypersyncgo

import (
	"github.com/enviodev/hypersync-client-go/pkg/options"
	"github.com/enviodev/hypersync-client-go/pkg/utils"
	"testing"
)

func TestClients(t *testing.T) {
	testCases := []struct {
		name string
		opts options.Node
	}{{
		name: "Test Ethereum Client",
		opts: options.Node{
			Type:      utils.EthereumNetwork,
			NetworkId: utils.EthereumNetworkID,
			Endpoint:  "https://eth.hypersync.xyz",
		},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Logf("Attempting to test %s client over %s", testCase.opts.Type, testCase.opts.GetEndpoint())
		})
	}
}
