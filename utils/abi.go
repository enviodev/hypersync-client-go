package utils

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"strings"

	abi "github.com/ethereum/go-ethereum/accounts/abi"
)

// MethodABI represents the JSON structure of a method's ABI.
type MethodABI struct {
	Name            string        `json:"name"`
	Type            string        `json:"type"` // "function", "constructor", etc.
	Inputs          []ArgumentABI `json:"inputs,omitempty"`
	Outputs         []ArgumentABI `json:"outputs,omitempty"`
	StateMutability string        `json:"stateMutability,omitempty"`
	Constant        bool          `json:"constant,omitempty"`
	Payable         bool          `json:"payable,omitempty"`
}

// ArgumentABI represents the JSON structure of a method's argument.
type ArgumentABI struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed,omitempty"`
}

// MethodToABI converts a Method object to its JSON ABI representation.
func MethodToABI(method *abi.Method) (string, error) {
	if method == nil {
		return "", fmt.Errorf("method is nil")
	}

	// Convert inputs and outputs
	inputs := make([]ArgumentABI, len(method.Inputs))
	for i, input := range method.Inputs {
		inputs[i] = ArgumentABI{
			Name:    input.Name,
			Type:    input.Type.String(),
			Indexed: input.Indexed,
		}
	}

	outputs := make([]ArgumentABI, len(method.Outputs))
	for i, output := range method.Outputs {
		outputs[i] = ArgumentABI{
			Name:    output.Name,
			Type:    output.Type.String(),
			Indexed: output.Indexed,
		}
	}

	methodABI := MethodABI{
		Name:            method.Name,
		Type:            "function",
		Inputs:          inputs,
		Outputs:         outputs,
		StateMutability: method.StateMutability,
		Constant:        method.Constant,
		Payable:         method.Payable,
	}

	abiJSON, err := json.Marshal([]MethodABI{methodABI})
	if err != nil {
		return "", fmt.Errorf("failed to marshal ABI: %s", err)
	}

	return string(abiJSON), nil
}

// EventToABI converts a Event object to its JSON ABI representation.
func EventToABI(method *abi.Event) (string, error) {
	if method == nil {
		return "", fmt.Errorf("event is nil")
	}

	// Convert inputs and outputs
	inputs := make([]ArgumentABI, len(method.Inputs))
	for i, input := range method.Inputs {
		inputs[i] = ArgumentABI{
			Name:    input.Name,
			Type:    input.Type.String(),
			Indexed: input.Indexed,
		}
	}

	// Construct the ABI entry
	methodABI := MethodABI{
		Name:   method.Name,
		Type:   "event",
		Inputs: inputs,
	}

	// Marshal to JSON
	abiJSON, err := json.Marshal([]MethodABI{methodABI})
	if err != nil {
		return "", fmt.Errorf("failed to marshal ABI: %s", err)
	}

	return string(abiJSON), nil
}

// ToABI converts the ABI object into an ethereum/go-ethereum ABI object.
func ToABI(data []byte) (*abi.ABI, error) {
	toReturn, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &toReturn, nil
}

// ToABIFromString converts the ABI string into an ethereum/go-ethereum ABI object.
func ToABIFromString(data string) (*abi.ABI, error) {
	toReturn, err := abi.JSON(strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &toReturn, nil
}
