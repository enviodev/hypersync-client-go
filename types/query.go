package types

type Query struct {
	FromBlock          uint64                 `json:"from_block,omitempty"`
	ToBlock            *uint64                `json:"to_block,omitempty"`
	Logs               []LogSelection         `json:"logs,omitempty"`
	Transactions       []TransactionSelection `json:"transactions,omitempty"`
	Traces             []TraceSelection       `json:"traces,omitempty"`
	IncludeAllBlocks   bool                   `json:"include_all_blocks,omitempty"`
	FieldSelection     FieldSelection         `json:"field_selection,omitempty"`
	MaxNumBlocks       *uint                  `json:"max_num_blocks,omitempty"`
	MaxNumTransactions *uint                  `json:"max_num_transactions,omitempty"`
	MaxNumLogs         *uint                  `json:"max_num_logs,omitempty"`
	MaxNumTraces       *uint                  `json:"max_num_traces,omitempty"`
	JoinMode           JoinMode               `json:"join_mode,omitempty"`
}

// QueryResponse represents the query response from hypersync instance.
// Contain next_block field in case query didn't process all the block range
type QueryResponse struct {
	// Current height of the source hypersync instance
	ArchiveHeight *int64 `json:"archive_height"`
	// Next block to query for, the responses are paginated so
	// the caller should continue the query from this block if they
	// didn't get responses up to the to_block they specified in the Query.
	NextBlock uint64 `json:"next_block"`
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

func (qr *QueryResponse) SetArchiveHeight(height *int64) {
	qr.ArchiveHeight = height
}

func (qr *QueryResponse) SetNextBlock(block uint64) {
	qr.NextBlock = block
}

func (qr *QueryResponse) SetTotalExecutionTime(tet uint64) {
	qr.TotalExecutionTime = tet
}

func (qr *QueryResponse) SetRollbackGuard(rg *RollbackGuard) {
	qr.RollbackGuard = rg
}
