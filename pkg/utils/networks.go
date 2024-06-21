package utils

import (
	"math/big"
)

// Network represents the network as a string.
type Network string

// NetworkID represents the network ID as an unsigned 64-bit integer.
type NetworkID uint64

// Predefined network IDs.
var (
	// EthereumNetworkID is the ID for the Ethereum network.
	EthereumNetworkID NetworkID = 1

	// EthereumNetwork is the string representation of the Ethereum network.
	EthereumNetwork Network = "ethereum"
)

// String returns the string representation of the NetworkID.
func (n NetworkID) String() string {
	return n.ToBig().String()
}

// IsValid checks if the NetworkID is valid (non-zero).
func (n NetworkID) IsValid() bool {
	return n != 0
}

// ToNetwork converts a NetworkID to its corresponding Network.
func (n NetworkID) ToNetwork() Network {
	switch n {
	case EthereumNetworkID:
		return EthereumNetwork
	default:
		return EthereumNetwork
	}
}

// Uint64 converts a NetworkID to a uint64.
func (n NetworkID) Uint64() uint64 {
	return uint64(n)
}

// ToBig converts a NetworkID to a big.Int.
func (n NetworkID) ToBig() *big.Int {
	return new(big.Int).SetUint64(uint64(n))
}
