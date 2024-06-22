package types

type TraceSelection struct {
	From       []Address `json:"from"`
	To         []Address `json:"to"`
	Address    []Address `json:"address"`
	CallType   []string  `json:"call_type"`
	RewardType []string  `json:"reward_type"`
	Kind       []string  `json:"type"`
	SigHash    []SigHash `json:"sighash"`
}

// Trace represents an Ethereum trace object.
type Trace struct {
	// The address of the sender who initiated the transaction.
	From *Address `json:"from,omitempty"`
	// The address of the recipient of the transaction if it was a transaction to an address. For contract creation transactions, this field is None.
	To *Address `json:"to,omitempty"`
	// The type of trace, `call` or `delegatecall`, two ways to invoke a function in a smart contract. `call` creates a new environment for the function to work in, so changes made in that function won't affect the environment where the function was called. `delegatecall` doesn't create a new environment. Instead, it runs the function within the environment of the caller, so changes made in that function will affect the caller's environment.
	CallType *string `json:"call_type,omitempty"`
	// The units of gas included in the transaction by the sender.
	Gas *Quantity `json:"gas,omitempty"`
	// The optional input data sent with the transaction, usually used to interact with smart contracts.
	Input *Data `json:"input,omitempty"`
	// The init code.
	Init *Data `json:"init,omitempty"`
	// The value of the native token transferred along with the transaction, in Wei.
	Value *Quantity `json:"value,omitempty"`
	// The address of the receiver for reward transaction.
	Author *Address `json:"author,omitempty"`
	// Kind of reward. `Block` reward or `Uncle` reward.
	RewardType *string `json:"reward_type,omitempty"`
	// The hash of the block in which the transaction was included.
	BlockHash *Hash `json:"block_hash,omitempty"`
	// The number of the block in which the transaction was included.
	BlockNumber *uint64 `json:"block_number,omitempty"`
	// Destroyed address.
	AddressDestroyed *Address `json:"address,omitempty"`
	// Contract code.
	Code *Data `json:"code,omitempty"`
	// The total used gas by the call, encoded as hexadecimal.
	GasUsed *Quantity `json:"gas_used,omitempty"`
	// The return value of the call, encoded as a hexadecimal string.
	Output *Data `json:"output,omitempty"`
	// The number of sub-traces created during execution. When a transaction is executed on the EVM, it may trigger additional sub-executions, such as when a smart contract calls another smart contract or when an external account is accessed.
	Subtraces *uint64 `json:"subtraces,omitempty"`
	// An array that indicates the position of the transaction in the trace.
	TraceAddress *[]uint64 `json:"trace_address,omitempty"`
	// The hash of the transaction.
	TransactionHash *Hash `json:"transaction_hash,omitempty"`
	// The index of the transaction in the block.
	TransactionPosition *uint64 `json:"transaction_position,omitempty"`
	// The type of action taken by the transaction, `call`, `create`, `reward` and `suicide`. `call` is the most common type of trace and occurs when a smart contract invokes another contract's function. `create` represents the creation of a new smart contract. This type of trace occurs when a smart contract is deployed to the blockchain.
	Kind *string `json:"type,omitempty"`
	// A string that indicates whether the transaction was successful or not. None if successful, Reverted if not.
	Error *string `json:"error,omitempty"`
}
