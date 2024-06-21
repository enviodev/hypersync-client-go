package hypersyncgo

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Structs to represent the JSON structure
type Response struct {
	Data               []Data        `json:"data"`
	ArchiveHeight      int           `json:"archive_height"`
	NextBlock          int           `json:"next_block"`
	TotalExecutionTime int           `json:"total_execution_time"`
	RollbackGuard      RollbackGuard `json:"rollback_guard"`
}

type RollbackGuard struct {
	BlockNumber      int    `json:"block_number"`
	Timestamp        int64  `json:"timestamp"`
	Hash             string `json:"hash"`
	FirstBlockNumber int    `json:"first_block_number"`
	FirstParentHash  string `json:"first_parent_hash"`
}

type Transaction struct {
	BlockNumber *big.Int    `json:"block_number"`
	Hash        common.Hash `json:"hash"`
}

type Data struct {
	Transactions []*Transaction `json:"transactions"`
}
