package types

// DataResponse represents the response containing various types of data
// such as blocks, transactions, logs, and traces.
type DataResponse struct {
	// Blocks contains the list of block data.
	Blocks []Block `json:"blocks,omitempty"`

	// Transactions contains the list of transaction data.
	Transactions []Transaction `json:"transactions,omitempty"`

	// Logs contains the list of log data.
	Logs []Log `json:"logs,omitempty"`

	// Traces contains the list of trace data.
	Traces []Trace `json:"traces,omitempty"`
}
