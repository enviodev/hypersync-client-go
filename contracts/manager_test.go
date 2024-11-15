package contracts

import (
	"github.com/enviodev/hypersync-client-go/utils"
	"testing"

	"github.com/enviodev/hypersync-client-go/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBasicManager(t *testing.T) {
	// Initialize the logger
	zLog, err := logger.GetZapDevelopmentLogger(zap.NewAtomicLevelAt(zap.DebugLevel))
	require.NoError(t, err)
	require.NotNil(t, zLog)
	logger.SetGlobalLogger(zLog)

	contract := defaultContracts[utils.EthereumNetworkID][0]
	contractAbi := contract.RawABI()

	// Define table-driven test cases
	tests := []struct {
		id           uint64
		name         string
		networkID    utils.NetworkID
		addr         string
		contractName string
		standard     utils.Standard
		abiRaw       string
		expectErr    bool
	}{
		{
			id:           1,
			name:         "Valid contract registration",
			networkID:    utils.NetworkID(1),
			addr:         "0x1234567890abcdef1234567890abcdef12345678",
			contractName: "ValidContract",
			standard:     utils.Erc20,
			abiRaw:       contractAbi,
			expectErr:    false,
		},
		{
			id:           1,
			name:         "Duplicate contract registration",
			networkID:    utils.NetworkID(1),
			addr:         "0x1234567890abcdef1234567890abcdef12345678",
			contractName: "DuplicateContract",
			standard:     utils.Erc20,
			abiRaw:       contractAbi,
			expectErr:    true, // Since the contract is already registered
		},
		{
			id:           3,
			name:         "Invalid contract ABI",
			networkID:    utils.NetworkID(3),
			addr:         "0x4567890abcdef1234567890abcdef1234567890",
			contractName: "InvalidABI",
			standard:     utils.NoStandard,
			abiRaw:       "",
			expectErr:    true,
		},
	}

	// Create a manager instance
	manager := NewManager()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a contract using the NewContract function
			addr := common.HexToAddress(tt.addr)
			contract, err := NewContract(tt.id, addr, tt.contractName, tt.standard, tt.abiRaw)

			// Check for contract creation error if expected
			if tt.expectErr && err != nil {
				require.Error(t, err, "expected error during contract creation")
				return
			} else {
				require.NoError(t, err, "unexpected error during contract creation")
			}

			// Register the contract
			err = manager.Register(tt.networkID, contract)
			if tt.expectErr {
				require.Error(t, err, "expected error during contract registration")
			} else {
				require.NoError(t, err, "unexpected error during contract registration")
			}
		})
	}
}

func TestNewManagerWithDefaults(t *testing.T) {
	// Initialize the logger
	zLog, err := logger.GetZapDevelopmentLogger(zap.NewAtomicLevelAt(zap.DebugLevel))
	require.NoError(t, err)
	require.NotNil(t, zLog)
	logger.SetGlobalLogger(zLog)

	// Define table-driven test cases
	tests := []struct {
		name      string
		setupFn   func() map[utils.NetworkID][]*Contract
		expectErr bool
	}{
		{
			name: "Successful manager creation with defaults",
			setupFn: func() map[utils.NetworkID][]*Contract {
				// Set up a valid default contracts map
				return defaultContracts
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the contracts to use in defaults
			dContracts := tt.setupFn()

			// Create the manager with defaults
			manager, err := NewManagerWithDefaults()

			// Check for error during manager creation
			if tt.expectErr {
				require.Error(t, err, "expected error during loading defaults")
				require.Nil(t, manager, "manager should be nil on error")
			} else {
				require.NoError(t, err, "unexpected error during manager creation")
				require.NotNil(t, manager, "manager should not be nil")

				// Check if all the expected default contracts were loaded correctly
				for networkID, expectedContracts := range dContracts {
					for _, expectedContract := range expectedContracts {
						loadedContract, err := manager.GetByID(networkID, expectedContract.ID())
						require.NoError(t, err)

						// Verify contract exists and matches the expected attributes
						require.NotNil(t, loadedContract, "expected contract to be present")
						require.Equal(t, expectedContract.Addr().Hex(), loadedContract.Addr().Hex(), "contract address should match")
						require.Equal(t, expectedContract.Name(), loadedContract.Name(), "contract name should match")
						require.Equal(t, expectedContract.Standard(), loadedContract.Standard(), "contract standard should match")
					}
				}
			}
		})
	}
}
