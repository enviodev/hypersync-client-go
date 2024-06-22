package types

import (
	"github.com/apache/arrow/go/v10/arrow"
)

type SigHash [4]byte
type Address string
type Hash string
type LogArgument string

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
	BlockNumber      uint64 `json:"block_number"`
	Timestamp        int64  `json:"timestamp"`
	Hash             Hash   `json:"hash"`
	FirstBlockNumber uint64 `json:"first_block_number"`
	FirstParentHash  Hash   `json:"first_parent_hash"`
}
