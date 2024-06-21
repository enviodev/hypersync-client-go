package types

type Query struct {
	FromBlock          uint64                 `json:"from_block"`
	ToBlock            *uint64                `json:"to_block"`
	Logs               []LogSelection         `json:"logs"`
	Transactions       []TransactionSelection `json:"transactions"`
	Traces             []TraceSelection       `json:"traces"`
	IncludeAllBlocks   bool                   `json:"include_all_blocks"`
	FieldSelection     FieldSelection         `json:"field_selection"`
	MaxNumBlocks       *uint                  `json:"max_num_blocks"`
	MaxNumTransactions *uint                  `json:"max_num_transactions"`
	MaxNumLogs         *uint                  `json:"max_num_logs"`
	MaxNumTraces       *uint                  `json:"max_num_traces"`
	JoinMode           JoinMode               `json:"join_mode"`
}
