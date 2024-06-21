package hypersyncgo

import "github.com/enviodev/hypersync-client-go/pkg/utils"

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
		if node.NetworkId == int(networkId.Uint64()) {
			return &node, true
		}
	}
	return nil, false
}

// Node represents the configuration and details of a network node.
type Node struct {
	// Group represents the group name of the node.
	Group string `mapstructure:"group" yaml:"group" json:"group"`

	// Type represents the type of the node.
	Type string `mapstructure:"type" yaml:"type" json:"type"`

	// NetworkId represents the network ID of the node.
	NetworkId int `mapstructure:"networkId" yaml:"networkId" json:"networkId"`

	// Endpoint represents the network endpoint of the node.
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint" json:"endpoint"`
}

// GetGroup returns the group name of the node.
func (n *Node) GetGroup() string {
	return n.Group
}

// GetType returns the type of the node.
func (n *Node) GetType() string {
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
