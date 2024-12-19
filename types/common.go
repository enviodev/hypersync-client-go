package types

import (
	"fmt"
	"math/big"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/ethereum/go-ethereum/common"
)

type DataType uint8

const (
	BlocksDataType DataType = iota
	TransactionsDataType
	LogsDataType
	TracesDataType
	// You can add more DataTypes here
)

type SigHash [4]byte

func NewSigHashFromHex(hex string) SigHash {
	return SigHash(common.Hex2Bytes(hex))
}

func (s SigHash) MarshalJSON() ([]byte, error) {
	return []byte(`"0x` + common.Bytes2Hex(s[:]) + `"`), nil
}

func (s *SigHash) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid hex string")
	}
	hexStr := string(data[1 : len(data)-1])
	bytes := common.Hex2Bytes(hexStr)
	if len(bytes) != 4 {
		return fmt.Errorf("invalid sighash length: expected 4 bytes, got %d", len(bytes))
	}
	copy(s[:], bytes)
	return nil
}

type Address string
type Hash string
type LogArgument string

type Nonce string
type BloomFilter string
type Quantity string
type Data string
type BlockNumber uint64
type TransactionIndex uint64
type LogIndex uint64
type Withdrawal struct{}
type TransactionType string
type TransactionStatus string

type JoinMode string

func hashDT() arrow.DataType {
	return &arrow.FixedSizeBinaryType{ByteWidth: 32}
}

func addrDT() arrow.DataType {
	return &arrow.FixedSizeBinaryType{ByteWidth: 20}
}

const (
	Default     JoinMode = "Default"
	JoinAll     JoinMode = "JoinAll"
	JoinNothing JoinMode = "JoinNothing"
)

type RollbackGuard struct {
	BlockNumber      *big.Int    `json:"block_number"`
	Timestamp        int64       `json:"timestamp"`
	Hash             common.Hash `json:"hash"`
	FirstBlockNumber uint64      `json:"first_block_number"`
	FirstParentHash  common.Hash `json:"first_parent_hash"`
}

// AccessList represents an Evm access list object.
//
// See ethereum rpc spec for the meaning of fields.
type AccessList struct {
	Address     *Address `json:"address,omitempty"`
	StorageKeys *[]Hash  `json:"storageKeys,omitempty"`
}
