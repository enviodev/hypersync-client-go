package types

import (
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

type LogSelection struct {
	Address []common.Address `json:"address,omitempty"`
	Topics  [][]common.Hash  `json:"topics,omitempty"`
}

// Log represents an Ethereum event log object.
type Log struct {
	// The boolean value indicating if the event was removed from the blockchain due to a chain reorganization. True if the log was removed. False if it is a valid log.
	Removed *bool `json:"removed,omitempty"`
	// The integer identifying the index of the event within the block's list of events.
	LogIndex *uint64 `json:"log_index,omitempty"`
	// The integer index of the transaction within the block's list of transactions.
	TransactionIndex *uint64 `json:"transaction_index,omitempty"`
	// The hash of the transaction that triggered the event.
	TransactionHash *common.Hash `json:"transaction_hash,omitempty"`
	// The hash of the block in which the event was included.
	BlockHash *common.Hash `json:"block_hash,omitempty"`
	// The block number in which the event was included.
	BlockNumber *big.Int `json:"block_number,omitempty"`
	// The contract address from which the event originated.
	Address *common.Address `json:"address,omitempty"`
	// The non-indexed data that was emitted along with the event.
	Data *[]byte `json:"data,omitempty"`
	// Additional topics
	Topic0 *common.Hash `json:"topic0,omitempty"`
	Topic1 *common.Hash `json:"topic1,omitempty"`
	Topic2 *common.Hash `json:"topic2,omitempty"`
	Topic3 *common.Hash `json:"topic3,omitempty"`
}

func NewLogFromRecord(schema *arrow.Schema, record arrow.Record) (*Log, error) {
	if record.NumCols() != int64(len(schema.Fields())) {
		return nil, errors.New("number of columns in record does not match schema")
	}

	toReturn := &Log{}

	for i, field := range schema.Fields() {
		col := record.Column(i)
		if col.Len() == 0 {
			continue
		}
		switch field.Name {
		case "removed":
			if fCol, ok := col.(*array.Boolean); ok {
				val := fCol.Value(0)
				toReturn.Removed = &val
			}
		case "log_index":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.LogIndex = &val
			}
		case "transaction_index":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.TransactionIndex = &val
			}
		case "transaction_hash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				hash := common.BytesToHash(val)
				toReturn.TransactionHash = &hash
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
				toReturn.Address = &address
			}
		case "data":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				toReturn.Data = &val
			}
		case "topic0":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				topic := common.BytesToHash(val)
				toReturn.Topic0 = &topic
			}
		case "topic1":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				topic := common.BytesToHash(val)
				toReturn.Topic1 = &topic
			}
		case "topic2":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				topic := common.BytesToHash(val)
				toReturn.Topic2 = &topic
			}
		case "topic3":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				topic := common.BytesToHash(val)
				toReturn.Topic3 = &topic
			}
		default:
			return nil, errors.New("unsupported field: " + field.Name)
		}
	}

	return toReturn, nil
}
