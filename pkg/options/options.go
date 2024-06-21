package options

import (
	"github.com/enviodev/hypersync-client-go/pkg/utils"
	"time"
)

// Options represents the configuration options for network nodes.
type Options struct {
	// Nodes is a slice of Node representing the network nodes.
	Blockchains []Node `mapstructure:"blockchains" yaml:"blockchains" json:"blockchains"`
}

// GetBlockchains returns the slice of network blockchains nodes from the Options.
func (o *Options) GetBlockchains() []Node {
	return o.Blockchains
}

func (o *Options) GetNodeByNetworkId(networkId utils.NetworkID) (*Node, bool) {
	for _, node := range o.Blockchains {
		if node.NetworkId == networkId {
			return &node, true
		}
	}
	return nil, false
}

// Node represents the configuration and details of a network node.
type Node struct {
	// Type represents the type of the node.
	Type utils.Network `mapstructure:"type" yaml:"type" json:"type"`

	// NetworkId represents the network ID of the node.
	NetworkId utils.NetworkID `mapstructure:"networkId" yaml:"networkId" json:"networkId"`

	// Endpoint represents the network endpoint of the node.
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint" json:"endpoint"`

	// BearerToken is the HyperSync server bearer token.
	BearerToken *string `mapstructure:"bearerToken" yaml:"bearerToken" json:"bearerToken"`

	// HTTPReqTimeoutMillis is the number of milliseconds to wait for a response before timing out.
	HTTPReqTimeoutMs *time.Duration `mapstructure:"httpReqTimeoutMs" yaml:"httpReqTimeoutMs" json:"httpReqTimeoutMs"`

	// MaxNumRetries is the number of retries to attempt before returning an error.
	MaxNumRetries int `mapstructure:"maxNumRetries" yaml:"maxNumRetries" json:"maxNumRetries" default:"12"`

	// RetryBackoffMs is the number of milliseconds used for retry backoff increasing.
	RetryBackoffMs time.Duration `mapstructure:"retryBackoffMs" yaml:"retryBackoffMs" json:"retryBackoffMs"`

	// RetryBaseMs is the initial wait time for request backoff.
	RetryBaseMs time.Duration `mapstructure:"retryBaseMs" yaml:"retryBaseMs" json:"retryBaseMs"`

	// RetryCeilingMs is the ceiling time for request backoff.
	RetryCeilingMs time.Duration `mapstructure:"retryCeilingMs" yaml:"retryCeilingMs" json:"retryCeilingMs"`
}

// GetType returns the type of the node.
func (n *Node) GetType() utils.Network {
	return n.Type
}

// GetNetworkID returns the network ID of the node.
func (n *Node) GetNetworkID() int64 {
	return int64(n.NetworkId)
}

// GetEndpoint returns the network endpoint of the node.
func (n *Node) GetEndpoint() string {
	return n.Endpoint
}
