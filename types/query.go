package types

import (
	"math/big"
)

type Query struct {
	FromBlock          *big.Int               `json:"from_block,omitempty"`
	ToBlock            *big.Int               `json:"to_block,omitempty"`
	Logs               []LogSelection         `json:"logs,omitempty"`
	Transactions       []TransactionSelection `json:"transactions,omitempty"`
	Traces             []TraceSelection       `json:"traces,omitempty"`
	IncludeAllBlocks   bool                   `json:"include_all_blocks,omitempty"`
	FieldSelection     FieldSelection         `json:"field_selection,omitempty"`
	MaxNumBlocks       *big.Int               `json:"max_num_blocks,omitempty"`
	MaxNumTransactions *big.Int               `json:"max_num_transactions,omitempty"`
	MaxNumLogs         *big.Int               `json:"max_num_logs,omitempty"`
	MaxNumTraces       *big.Int               `json:"max_num_traces,omitempty"`
	JoinMode           JoinMode               `json:"join_mode,omitempty"`
}

// QueryResponse represents the query response from hypersync instance.
// Contain next_block field in case query didn't process all the block range
type QueryResponse struct {
	// Current height of the source hypersync instance
	ArchiveHeight *big.Int `json:"archive_height"`
	// Next block to query for, the responses are paginated so
	// the caller should continue the query from this block if they
	// didn't get responses up to the to_block they specified in the Query.
	NextBlock *big.Int `json:"next_block"`
	// Total time it took the hypersync instance to execute the query.
	TotalExecutionTime uint64 `json:"total_execution_time"`
	// Response data
	Data DataResponse `json:"data"`
	// Rollback guard
	RollbackGuard *RollbackGuard `json:"rollback_guard"`
}

func (qr *QueryResponse) GetData() DataResponse {
	return qr.Data
}

func (qr *QueryResponse) AppendBlockData(data Block) {
	qr.Data.Blocks = append(qr.Data.Blocks, data)
}

func (qr *QueryResponse) AppendTransactionData(data Transaction) {
	qr.Data.Transactions = append(qr.Data.Transactions, data)
}

func (qr *QueryResponse) AppendLogData(data Log) {
	qr.Data.Logs = append(qr.Data.Logs, data)
}

func (qr *QueryResponse) AppendTraceData(data Trace) {
	qr.Data.Traces = append(qr.Data.Traces, data)
}

func (qr *QueryResponse) SetArchiveHeight(height *big.Int) {
	qr.ArchiveHeight = height
}

func (qr *QueryResponse) SetNextBlock(block *big.Int) {
	qr.NextBlock = block
}

func (qr *QueryResponse) HasNextBlock() bool {
	return qr.NextBlock != nil
}

func (qr *QueryResponse) SetTotalExecutionTime(tet uint64) {
	qr.TotalExecutionTime = tet
}

func (qr *QueryResponse) SetRollbackGuard(rg *RollbackGuard) {
	qr.RollbackGuard = rg
}

func (qr *QueryResponse) HasRollbackGuard() bool {
	return qr.RollbackGuard != nil
}

func (qr *QueryResponse) GetRollbackGuard() *RollbackGuard {
	return qr.RollbackGuard
}

func (qr *QueryResponse) GetBlocks() []Block {
	return qr.Data.Blocks
}

func (qr *QueryResponse) GetBlockByNumber(number *big.Int) *Block {
	for _, block := range qr.Data.Blocks {
		if block.Number.Cmp(number) == 0 {
			return &block
		}
	}
	return nil
}

func (qr *QueryResponse) GetTransactions() []Transaction {
	return qr.Data.Transactions
}

func (qr *QueryResponse) GetLogs() []Log {
	return qr.Data.Logs
}

func (qr *QueryResponse) GetTraces() []Trace {
	return qr.Data.Traces
}
