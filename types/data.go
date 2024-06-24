package types

type DataResponse struct {
	Blocks       []Block       `json:"blocks,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
	Logs         []Log         `json:"logs,omitempty"`
	Traces       []Trace       `json:"traces,omitempty"`
}
