package types

type LogSelection struct {
	Address []Address        `json:"address"`
	Topics  [4][]LogArgument `json:"topics"`
}

// Log represents an Ethereum event log object.
type Log struct {
	// The boolean value indicating if the event was removed from the blockchain due to a chain reorganization. True if the log was removed. False if it is a valid log.
	Removed *bool `json:"removed,omitempty"`
	// The integer identifying the index of the event within the block's list of events.
	LogIndex *LogIndex `json:"log_index,omitempty"`
	// The integer index of the transaction within the block's list of transactions.
	TransactionIndex *TransactionIndex `json:"transaction_index,omitempty"`
	// The hash of the transaction that triggered the event.
	TransactionHash *Hash `json:"transaction_hash,omitempty"`
	// The hash of the block in which the event was included.
	BlockHash *Hash `json:"block_hash,omitempty"`
	// The block number in which the event was included.
	BlockNumber *BlockNumber `json:"block_number,omitempty"`
	// The contract address from which the event originated.
	Address *Address `json:"address,omitempty"`
	// The non-indexed data that was emitted along with the event.
	Data *Data `json:"data,omitempty"`
	// An array of 32-byte data fields containing indexed event parameters.
	Topics [4]*LogArgument `json:"topics"`
}
