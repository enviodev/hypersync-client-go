package hypersyncgo

import (
	"context"
	"encoding/json"
	"github.com/enviodev/hypersync-client-go/options"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientRequiresApiToken(t *testing.T) {
	ctx := context.Background()

	_, err := NewClient(ctx, options.Node{
		Type:        utils.EthereumNetwork,
		NetworkId:   utils.EthereumNetworkID,
		Endpoint:    "https://eth.hypersync.xyz",
		RpcEndpoint: "https://eth.rpc.hypersync.xyz",
		ApiToken:    "",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "api token is required")
}

func TestNewClientRequiresEndpoint(t *testing.T) {
	ctx := context.Background()

	_, err := NewClient(ctx, options.Node{
		Type:        utils.EthereumNetwork,
		NetworkId:   utils.EthereumNetworkID,
		Endpoint:    "",
		RpcEndpoint: "https://eth.rpc.hypersync.xyz",
		ApiToken:    "my-token",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "endpoint is required")
}

func TestApiTokenSentOnRequests(t *testing.T) {
	expectedToken := "my-secret-api-token"

	// Create a test HTTP server that checks for the Authorization header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+expectedToken, authHeader, "Authorization header should contain the API token")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return a minimal valid JSON response based on the path
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/height":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"height": "0x1"})
		default:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		}
	}))
	defer server.Close()

	// Create a client that skips RPC by directly constructing the struct
	client := &Client{
		ctx: context.Background(),
		opts: options.Node{
			Type:        utils.EthereumNetwork,
			NetworkId:   utils.EthereumNetworkID,
			Endpoint:    server.URL,
			ApiToken:    expectedToken,
		},
		client: server.Client(),
	}

	// Test DoQuery sends the API token
	t.Run("DoQuery sends API token", func(t *testing.T) {
		query := &types.Query{
			FromBlock: big.NewInt(1),
			FieldSelection: types.FieldSelection{
				Block: []string{"number"},
			},
		}
		// DoQuery will get a response that may not parse perfectly, but we're
		// testing the header is sent, which is verified by the server handler above.
		_, _ = DoQuery[*types.Query, map[string]interface{}](context.Background(), client, http.MethodPost, query)
	})

	// Test Do sends the API token
	t.Run("Do sends API token", func(t *testing.T) {
		_, _ = Do[struct{}, map[string]interface{}](context.Background(), client, server.URL+"/height", http.MethodGet, struct{}{})
	})
}

func TestOptionsValidateRejectsMissingToken(t *testing.T) {
	opts := options.Options{
		Blockchains: []options.Node{
			{
				Type:        utils.EthereumNetwork,
				NetworkId:   utils.EthereumNetworkID,
				Endpoint:    "https://eth.hypersync.xyz",
				RpcEndpoint: "https://eth.rpc.hypersync.xyz",
				ApiToken:    "",
			},
		},
	}

	err := opts.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "api token is required")
}

func TestOptionsValidateAcceptsValidConfig(t *testing.T) {
	opts := options.Options{
		Blockchains: []options.Node{
			{
				Type:        utils.EthereumNetwork,
				NetworkId:   utils.EthereumNetworkID,
				Endpoint:    "https://eth.hypersync.xyz",
				RpcEndpoint: "https://eth.rpc.hypersync.xyz",
				ApiToken:    "valid-token",
			},
		},
	}

	err := opts.Validate()
	require.NoError(t, err)
}
