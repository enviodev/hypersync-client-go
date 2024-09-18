package decoder

import (
	"fmt"
	"github.com/enviodev/hypersync-client-go/contracts"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

// EthereumTopic represents a single decoded topic from an Ethereum event log. Topics are attributes
// of an event, such as the method signature and indexed parameters.
type EthereumTopic struct {
	Name  string `json:"name"`  // The name of the topic.
	Value any    `json:"value"` // The value of the topic, decoded into the appropriate Go data type.
}

// EthereumLog encapsulates a decoded Ethereum event log. It includes the event's details such as its name,
// signature, the contract that emitted the event, and the decoded data and topics.
type EthereumLog struct {
	Event        *abi.Event      `json:"-"`             // ABI definition of the log's event.
	Abi          string          `json:"abi"`           // ABI string of the event.
	SignatureHex common.Hash     `json:"signature_hex"` // Hex-encoded signature of the event.
	Signature    string          `json:"signature"`     // Signature of the event.
	Type         string          `json:"type"`          // Type of the event
	Name         string          `json:"name"`          // Name of the event.
	Data         map[string]any  `json:"data"`          // Decoded event data.
	Topics       []EthereumTopic `json:"topics"`        // Decoded topics of the event.
}

// DecodeEthereumLogWithContract decodes an Ethereum event log using the provided contract's ABI.
// It returns an EthereumLog instance containing the decoded event details such as the event name, data, and topics.
//
// Example:
//
//	log := &types.Log{}
//	contract := &contracts.Contract{}
//	decodedLog, err := DecodeEthereumLogWithContract(log, contract)
//	if err != nil {
//	    log.Fatalf("Failed to decode log: %v", err)
//	}
func DecodeEthereumLogWithContract(log types.Log, contract *contracts.Contract) (*EthereumLog, error) {
	return DecodeEthereumLog(log, contract.RawABI())
}

// DecodeEthereumLog decodes an Ethereum event log using the provided ABI data.
// It returns an EthereumLog instance containing the decoded event name, data, and topics.
//
// Example:
//
//	log := &types.Log{}
//	abiData := "<ABI JSON string>"
//	decodedLog, err := DecodeEthereumLog(log, abiData)
//	if err != nil {
//	    log.Fatalf("Failed to decode log: %v", err)
//	}
func DecodeEthereumLog(log types.Log, aData string) (*EthereumLog, error) {
	if len(log.Topics()) < 1 {
		return nil, errors.New("log is nil or has no topics")
	}

	topics := log.Topics()
	logABI, err := abi.JSON(strings.NewReader(aData))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse abi")
	}

	event, err := logABI.EventByID(topics[0])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get event by topic0: %s", topics[0].Hex())
	}

	data := make(map[string]any)
	if len(log.GetData()) > 0 {
		if uErr := event.Inputs.UnpackIntoMap(data, log.GetData()); uErr != nil {
			return nil, errors.Wrap(uErr, "failed to unpack inputs into map")
		}
	}

	// Identify and decode indexed inputs
	indexedInputs := make([]abi.Argument, 0)
	for _, input := range event.Inputs {
		if input.Indexed {
			indexedInputs = append(indexedInputs, input)
		}
	}

	decodedTopics := make([]EthereumTopic, len(indexedInputs))
	if len(topics) >= len(indexedInputs)+1 {
		for i, indexedInput := range indexedInputs {
			decodedTopic, dtErr := decodeTopic(topics[i+1], indexedInput)
			if dtErr != nil {
				return nil, errors.Wrap(dtErr, "failed to decode topic")
			}

			decodedTopics[i] = EthereumTopic{
				Name:  indexedInput.Name,
				Value: decodedTopic,
			}
		}
	}

	eventAbi, err := utils.EventToABI(event)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cast event into the abi")
	}

	toReturn := &EthereumLog{
		Event:        event,
		Abi:          eventAbi,
		SignatureHex: topics[0],
		Signature:    strings.TrimLeft(event.String(), "event "),
		Name:         event.Name,
		Type:         strings.ToLower(event.Name),
		Data:         data,
		Topics:       decodedTopics,
	}

	return toReturn, nil
}

// decodeTopic decodes a single topic from an Ethereum event log based on its ABI argument type.
// It supports various data types including addresses, booleans, integers, strings, bytes, and more.
//
// Example:
//
//	topic := common.Hash{}
//	argument := abi.Argument{}
//	decodedTopic, err := decodeTopic(topic, argument)
//	if err != nil {
//	    log.Fatalf("Failed to decode topic: %v", err)
//	}
func decodeTopic(topic common.Hash, argument abi.Argument) (interface{}, error) {
	switch argument.Type.T {
	case abi.AddressTy:
		return common.BytesToAddress(topic.Bytes()), nil
	case abi.BoolTy:
		return topic[common.HashLength-1] == 1, nil
	case abi.IntTy, abi.UintTy:
		size := argument.Type.Size
		if size > 256 {
			return nil, fmt.Errorf("unsupported integer size: %d", size)
		}
		integer := new(big.Int).SetBytes(topic[:])
		if argument.Type.T == abi.IntTy && size < 256 {
			integer = adjustIntSize(integer, size)
		}
		return integer, nil
	case abi.StringTy:
		return topic, nil
	case abi.BytesTy, abi.FixedBytesTy:
		return topic.Bytes(), nil
	case abi.SliceTy, abi.ArrayTy:
		return nil, fmt.Errorf("array/slice decoding not implemented")
	default:
		return nil, fmt.Errorf("decoding for type %v not implemented", argument.Type.T)
	}
}

// GetEthereumTopicByName searches for and returns a Topic by its name from a slice of Topic instances.
// It facilitates accessing specific topics directly by name rather than iterating over the slice.
// If the topic is not found, it returns nil.
func GetEthereumTopicByName(name string, topics []EthereumTopic) *EthereumTopic {
	for _, topic := range topics {
		if topic.Name == name {
			return &topic
		}
	}
	return nil
}

// adjustIntSize adjusts the size of an integer to match its ABI-specified size, which is relevant
// for signed integers smaller than 256 bits. This function ensures the integer is correctly
// interpreted according to its defined bit size in the ABI.
func adjustIntSize(integer *big.Int, size int) *big.Int {
	if size == 256 || integer.Bit(size-1) == 0 {
		return integer
	}
	return new(big.Int).Sub(integer, new(big.Int).Lsh(big.NewInt(1), uint(size)))
}
