package contracts

import (
	"fmt"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"strings"
)

// Contract represents an Ethereum smart contract, including its address, name, standard, and ABI.
type Contract struct {
	id       uint64
	addr     common.Address
	name     string
	standard utils.Standard
	rawAbi   string
	abi      *abi.ABI
}

// NewContract creates a new Contract instance, validates it, and returns an error if validation fails.
//
// Example:
//
//	addr := common.HexToAddress("0x123...")
//	contract, err := contracts.NewContract(addr, "MyContract", utils.Erc20, contractABI)
//	if err != nil {
//	    log.Fatalf("Failed to create contract: %v", err)
//	}
func NewContract(id uint64, addr common.Address, name string, standard utils.Standard, rawAbi string) (*Contract, error) {
	contract := &Contract{
		id:       id,
		addr:     addr,
		name:     name,
		standard: standard,
		rawAbi:   rawAbi,
	}

	if err := contract.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid contract")
	}

	return contract, nil
}

// Validate checks if the contract has a valid address, standard, and ABI, and parses the ABI.
//
// Example:
//
//	if err := contract.Validate(); err != nil {
//	    log.Fatalf("Invalid contract: %v", err)
//	}
func (c *Contract) Validate() error {
	if c.id == 0 {
		return errors.New("contract id is not set")
	}

	if !common.IsHexAddress(c.addr.Hex()) {
		return fmt.Errorf("invalid contract address provided: %s", c.addr.Hex())
	}

	if strings.TrimSpace(c.rawAbi) == "" {
		return errors.New("contract ABI is empty")
	}

	// Check if the contract standard is valid (e.g., ERC-20, ERC-721)
	if c.standard == utils.NoStandard {
		return errors.New("contract standard is not defined")
	}

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(c.rawAbi))
	if err != nil {
		return errors.Wrap(err, "invalid contract ABI")
	}

	c.abi = &parsedABI
	return nil
}

// ID returns the contract internal ID.
//
// Example:
//
//	fmt.Println(contract.ID())
func (c *Contract) ID() uint64 {
	return c.id
}

// Addr returns the Ethereum address of the contract.
//
// Example:
//
//	fmt.Println(contract.Addr().Hex())
func (c *Contract) Addr() common.Address {
	return c.addr
}

// Name returns the name of the contract.
//
// Example:
//
//	fmt.Println(contract.Name())
func (c *Contract) Name() string {
	return c.name
}

// Standard returns the standard of the contract (e.g., ERC-20, ERC-721).
//
// Example:
//
//	fmt.Println(contract.Standard())
func (c *Contract) Standard() utils.Standard {
	return c.standard
}

// RawABI returns the raw ABI string of the contract.
//
// Example:
//
//	fmt.Println(contract.RawAbi())
func (c *Contract) RawABI() string {
	return c.rawAbi
}

// ABI returns the parsed ABI of the contract.
//
// Example:
//
//	fmt.Println(contract.ABI())
func (c *Contract) ABI() *abi.ABI {
	return c.abi
}

// ToABI parses and returns the contract's ABI as an abi.ABI object, or an error if the ABI is invalid.
//
// Example:
//
//	parsedABI, err := contract.ToABI()
//	if err != nil {
//	    log.Fatalf("Failed to parse contract ABI: %v", err)
//	}
func (c *Contract) ToABI() (abi.ABI, error) {
	if c.abi == nil {
		// Parse the raw ABI if it's not already parsed
		parsedABI, err := abi.JSON(strings.NewReader(c.rawAbi))
		if err != nil {
			return abi.ABI{}, errors.Wrap(err, "invalid contract ABI")
		}
		c.abi = &parsedABI
	}
	return *c.abi, nil
}

// IsValid checks if the contract is valid and ready to be used.
// This is a helper method to check if the contract's key fields (addr, abi, and standard) are valid.
//
// Example:
//
//	if contract.IsValid() {
//	    // Use the contract
//	}
func (c *Contract) IsValid() bool {
	return c.addr != (common.Address{}) && c.abi != nil && c.standard != utils.NoStandard
}
