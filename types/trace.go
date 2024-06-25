package types

import (
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

type TraceSelection struct {
	From       []common.Address `json:"from,omitempty"`
	To         []common.Address `json:"to,omitempty"`
	Address    []common.Address `json:"address,omitempty"`
	CallType   []string         `json:"call_type,omitempty"`
	RewardType []string         `json:"reward_type,omitempty"`
	Kind       []string         `json:"type,omitempty"`
	SigHash    []common.Hash    `json:"sighash,omitempty"`
}

// Trace represents an Ethereum trace object.
type Trace struct {
	// The address of the sender who initiated the transaction.
	From *common.Address `json:"from,omitempty"`
	// The address of the recipient of the transaction if it was a transaction to an address. For contract creation transactions, this field is None.
	To *common.Address `json:"to,omitempty"`
	// SigHash ...
	SigHash *common.Hash `json:"sig_hash,omitempty"`
	// The type of trace, `call` or `delegatecall`, two ways to invoke a function in a smart contract. `call` creates a new environment for the function to work in, so changes made in that function won't affect the environment where the function was called. `delegatecall` doesn't create a new environment. Instead, it runs the function within the environment of the caller, so changes made in that function will affect the caller's environment.
	CallType *string `json:"call_type,omitempty"`
	// The units of gas included in the transaction by the sender.
	Gas *uint64 `json:"gas,omitempty"`
	// The optional input data sent with the transaction, usually used to interact with smart contracts.
	Input *[]byte `json:"input,omitempty"`
	// The init code.
	Init *[]common.Hash `json:"init,omitempty"`
	// The value of the native token transferred along with the transaction, in Wei.
	Value *big.Int `json:"value,omitempty"`
	// The address of the receiver for reward transaction.
	Author *common.Address `json:"author,omitempty"`
	// Kind of reward. `Block` reward or `Uncle` reward.
	RewardType *string `json:"reward_type,omitempty"`
	// The hash of the block in which the transaction was included.
	BlockHash *common.Hash `json:"block_hash,omitempty"`
	// The number of the block in which the transaction was included.
	BlockNumber *big.Int `json:"block_number,omitempty"`
	// Destroyed address.
	AddressDestroyed *common.Address `json:"address,omitempty"`
	// Contract code.
	Code *common.Hash `json:"code,omitempty"`
	// The total used gas by the call, encoded as hexadecimal.
	GasUsed *uint64 `json:"gas_used,omitempty"`
	// The return value of the call, encoded as a hexadecimal string.
	Output *common.Hash `json:"output,omitempty"`
	// The number of sub-traces created during execution. When a transaction is executed on the EVM, it may trigger additional sub-executions, such as when a smart contract calls another smart contract or when an external account is accessed.
	Subtraces *uint64 `json:"subtraces,omitempty"`
	// An array that indicates the position of the transaction in the trace.
	TraceAddress *[]uint64 `json:"trace_address,omitempty"`
	// The hash of the transaction.
	TransactionHash *common.Hash `json:"transaction_hash,omitempty"`
	// The index of the transaction in the block.
	TransactionPosition *uint64 `json:"transaction_position,omitempty"`
	// The type of action taken by the transaction, `call`, `create`, `reward` and `suicide`. `call` is the most common type of trace and occurs when a smart contract invokes another contract's function. `create` represents the creation of a new smart contract. This type of trace occurs when a smart contract is deployed to the blockchain.
	Kind *string `json:"type,omitempty"`
	// A string that indicates whether the transaction was successful or not. None if successful, Reverted if not.
	Error *string `json:"error,omitempty"`
}

func NewTraceFromRecord(schema *arrow.Schema, record arrow.Record) (*Trace, error) {
	if record.NumCols() != int64(len(schema.Fields())) {
		return nil, errors.New("number of columns in record does not match schema")
	}

	toReturn := &Trace{}

	for i, field := range schema.Fields() {
		col := record.Column(i)
		if col.Len() == 0 {
			continue
		}
		switch field.Name {
		case "from":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				from := common.BytesToAddress(val)
				toReturn.From = &from
			}
		case "to":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				to := common.BytesToAddress(val)
				toReturn.To = &to
			}
		case "sighash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				hash := common.BytesToHash(val)
				toReturn.SigHash = &hash
			}
		case "call_type":
			if fCol, ok := col.(*array.String); ok {
				val := fCol.Value(0)
				toReturn.CallType = &val
			}
		case "gas":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.Gas = &val
			}
		case "input":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				toReturn.Input = &val
			}
		case "init":
			if fCol, ok := col.(*array.List); ok {
				list := fCol.ListValues()
				init := make([]common.Hash, list.Len())
				/*				for j := 0; j < list.Len(); j++ {
								init[j] = common.BytesToHash(list.Value(j).([]byte))
							}*/
				toReturn.Init = &init
			}
		case "value":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.Value = big.NewInt(0).SetUint64(val)
			}
		case "author":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				author := common.BytesToAddress(val)
				toReturn.Author = &author
			}
		case "reward_type":
			if fCol, ok := col.(*array.String); ok {
				val := fCol.Value(0)
				toReturn.RewardType = &val
			}
		case "block_hash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				blockHash := common.BytesToHash(val)
				toReturn.BlockHash = &blockHash
			}
		case "block_number":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.BlockNumber = big.NewInt(0).SetUint64(val)
			}
		case "address":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				address := common.BytesToAddress(val)
				toReturn.AddressDestroyed = &address
			}
		case "code":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				code := common.BytesToHash(val)
				toReturn.Code = &code
			}
		case "gas_used":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.GasUsed = &val
			}
		case "output":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				output := common.BytesToHash(val)
				toReturn.Output = &output
			}
		case "subtraces":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.Subtraces = &val
			}
		case "trace_address":
			if fCol, ok := col.(*array.List); ok {
				list := fCol.ListValues()
				traceAddress := make([]uint64, list.Len())
				/*				for j := 0; j < list.Len(); j++ {
								traceAddress[j] = list.Value(j).(uint64)
							}*/
				toReturn.TraceAddress = &traceAddress
			}
		case "transaction_hash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				transactionHash := common.BytesToHash(val)
				toReturn.TransactionHash = &transactionHash
			}
		case "transaction_position":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.TransactionPosition = &val
			}
		case "type":
			if fCol, ok := col.(*array.String); ok {
				val := fCol.Value(0)
				toReturn.Kind = &val
			}
		case "error":
			if fCol, ok := col.(*array.String); ok {
				val := fCol.Value(0)
				toReturn.Error = &val
			}
		default:
			return nil, errors.New("unsupported field: " + field.Name)
		}
	}

	return toReturn, nil
}
