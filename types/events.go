package types

// Event represents an Ethereum event object.
type Event struct {
	// An Ethereum event transaction object.
	Transaction *Transaction `json:"transaction,omitempty"`
	// An Ethereum event block object.
	Block *Block `json:"block,omitempty"`
	// An Ethereum event log object.
	Log Log `json:"log"`
}
