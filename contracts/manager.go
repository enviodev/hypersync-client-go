package contracts

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"sync"

	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/pkg/errors"
)

// Manager manages a registry of contracts, allowing for thread-safe operations.
type Manager struct {
	mu       sync.RWMutex
	registry map[utils.NetworkID][]*Contract // Registry now holds a slice of contracts for each network.
}

// NewManager creates and returns a new Manager with an initialized registry.
func NewManager() *Manager {
	return &Manager{
		registry: make(map[utils.NetworkID][]*Contract),
	}
}

// NewManagerWithDefaults creates and returns a new Manager with an initialized registry and loaded default contracts.
func NewManagerWithDefaults() (*Manager, error) {
	m := NewManager()
	if err := m.LoadDefaults(); err != nil {
		return nil, err
	}
	return m, nil
}

// Register adds a new contract to the registry for the given network ID.
// It returns an error if the contract is already registered or if the contract is invalid.
//
// Example:
//
//	err := manager.Register(networkID, contract)
//	if err != nil {
//	    log.Fatalf("Failed to register contract: %v", err)
//	}
func (m *Manager) Register(network utils.NetworkID, contract *Contract) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the contract with the same address is already registered.
	for _, existingContract := range m.registry[network] {
		if existingContract.id == contract.id {
			return errors.Errorf("contract already registered with id: %d for network: %d", contract.id, network)
		}
	}

	// Validate the contract before adding it to the registry.
	if err := contract.Validate(); err != nil {
		return errors.Wrap(err, "contract validation failed")
	}

	// Append the contract to the list for the given network ID.
	m.registry[network] = append(m.registry[network], contract)
	return nil
}

// GetByAddr retrieves the contract for the given network ID and associated address.
// It returns an error if no contract is registered for the given network.
//
// Example:
//
//	contract, err := manager.GetByAddr(networkID, common.HexToAddress("0x0"))
//	if err != nil {
//	    log.Fatalf("Failed to get contract: %v", err)
//	}
func (m *Manager) GetByAddr(network utils.NetworkID, address common.Address) (*Contract, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	contracts, exists := m.registry[network]
	if !exists || len(contracts) == 0 {
		return nil, errors.Errorf("no contracts registered for network: %d", network)
	}

	for _, contract := range contracts {
		if contract.addr == address {
			return contract, nil
		}
	}

	return nil, fmt.Errorf("requested contract could not be found for network: %d, address: %s", network, address.Hex())
}

// GetByID retrieves the contract for the given network ID and associated contract ID.
// It returns an error if no contract is registered for the given network.
//
// Example:
//
//	contract, err := manager.GetByAddr(networkID, 1)
//	if err != nil {
//	    log.Fatalf("Failed to get contract: %v", err)
//	}
func (m *Manager) GetByID(network utils.NetworkID, id uint64) (*Contract, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	contracts, exists := m.registry[network]
	if !exists || len(contracts) == 0 {
		return nil, errors.Errorf("no contracts registered for network: %d", network)
	}

	for _, contract := range contracts {
		if contract.id == id {
			return contract, nil
		}
	}

	return nil, fmt.Errorf("requested contract could not be found for network: %d, id: %d", network, id)
}

// List returns a copy of all registered contracts in the registry.
//
// Example:
//
//	contracts := manager.List()
//	for id, contract := range contracts {
//	    fmt.Printf("Network: %d, Contract: %v\n", id, contract)
//	}
func (m *Manager) List() map[utils.NetworkID][]*Contract {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy of the map to avoid potential concurrent access issues.
	toReturn := make(map[utils.NetworkID][]*Contract)
	for id, contract := range m.registry {
		toReturn[id] = contract
	}

	return toReturn
}

// ListByNetworkId returns a copy of all registered contracts in the registry under specified network id.
//
// Example:
//
//	contracts := manager.ListByNetworkId(utils.EthereumNetworkId)
//	for id, contract := range contracts {
//	    fmt.Printf("Network: %d, Contract: %v\n", id, contract)
//	}
func (m *Manager) ListByNetworkId(networkId utils.NetworkID) []*Contract {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy of the map to avoid potential concurrent access issues.
	toReturn := make([]*Contract, 0)
	for _, contract := range m.registry {
		toReturn = append(toReturn, contract...)
	}

	return toReturn
}
