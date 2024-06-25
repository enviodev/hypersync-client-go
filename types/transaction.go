package types

import (
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"math/big"
)

type TransactionSelection struct {
	From            []common.Address `json:"from,omitempty"`
	To              []common.Address `json:"to,omitempty"`
	SigHash         []SigHash        `json:"sighash,omitempty"`
	Status          *uint8           `json:"status,omitempty"`
	Kind            []uint8          `json:"type,omitempty"`
	ContractAddress []common.Address `json:"contract_address,omitempty"`
}

// Transaction represents an Ethereum transaction object.
type Transaction struct {
	// The Keccak 256-bit hash of the block
	BlockHash *common.Hash `json:"block_hash,omitempty"`
	// A scalar value equal to the number of ancestor blocks. The genesis block has a number of zero; formally Hi.
	BlockNumber *big.Int `json:"block_number,omitempty"`
	// SigHash ...
	SigHash *common.Hash `json:"sig_hash,omitempty"`
	// The 160-bit address of the message call’s sender
	From *common.Address `json:"from,omitempty"`
	// A scalar value equal to the maximum amount of gas that should be used in executing this transaction. This is paid up-front, before any computation is done and may not be increased later; formally Tg.
	Gas *uint64 `json:"gas,omitempty"`
	// A scalar value equal to the number of Wei to be paid per unit of gas for all computation costs incurred as a result of the execution of this transaction; formally Tp.
	GasPrice *big.Int `json:"gas_price,omitempty"`
	// A transaction hash is a keccak hash of an RLP encoded signed transaction.
	Hash *common.Hash `json:"hash,omitempty"`
	// Input has two uses depending if transaction is Create or Call (if `to` field is None or Some). pub init: An unlimited size byte array specifying the EVM-code for the account initialisation procedure CREATE, data: An unlimited size byte array specifying the input data of the message call, formally Td.
	Input *[]byte `json:"input,omitempty"`
	// A scalar value equal to the number of transactions sent by the sender; formally Tn.
	Nonce *uint64 `json:"nonce,omitempty"`
	// The 160-bit address of the message call’s recipient or, for a contract creation transaction, ∅, used here to denote the only member of B0 ; formally Tt.
	To *common.Address `json:"to,omitempty"`
	// Index of the transaction in the block
	TransactionIndex *uint64 `json:"transaction_index,omitempty"`
	// A scalar value equal to the number of Wei to be transferred to the message call’s recipient or, in the case of contract creation, as an endowment to the newly created account; formally Tv.
	Value *big.Int `json:"value,omitempty"`
	// Replay protection value based on chain_id. See EIP-155 for more info.
	V *big.Int `json:"v,omitempty"`
	// The R field of the signature; the point on the curve.
	R *big.Int `json:"r,omitempty"`
	// The S field of the signature; the point on the curve.
	S *big.Int `json:"s,omitempty"`
	// yParity: Signature Y parity; formally Ty
	YParity *big.Int `json:"y_parity,omitempty"`
	// Max Priority fee that transaction is paying As ethereum circulation is around 120mil eth as of 2022 that is around 120000000000000000000000000 wei we are safe to use u128 as its max number is: 340282366920938463463374607431768211455 This is also known as `GasTipCap`
	MaxPriorityFeePerGas *big.Int `json:"max_priority_fee_per_gas,omitempty"`
	// A scalar value equal to the maximum amount of gas that should be used in executing this transaction. This is paid up-front, before any computation is done and may not be increased later; formally Tg. As ethereum circulation is around 120mil eth as of 2022 that is around 120000000000000000000000000 wei we are safe to use u128 as its max number is: 340282366920938463463374607431768211455 This is also known as `GasFeeCap`
	MaxFeePerGas *big.Int `json:"max_fee_per_gas,omitempty"`
	// Added as EIP-155: Simple replay attack protection
	ChainID *big.Int `json:"chain_id,omitempty"`
	// The accessList specifies a list of addresses and storage keys; these addresses and storage keys are added into the `accessed_addresses` and `accessed_storage_keys` global sets (introduced in EIP-2929). A gas cost is charged, though at a discount relative to the cost of accessing outside the list.
	AccessList *[]types.AccessList `json:"access_list,omitempty"`
	// Max fee per data gas aka BlobFeeCap or blobGasFeeCap
	MaxFeePerBlobGas *big.Int `json:"max_fee_per_blob_gas,omitempty"`
	// It contains a vector of fixed size hash(32 bytes)
	BlobVersionedHashes *[]common.Hash `json:"blob_versioned_hashes,omitempty"`
	// The total amount of gas used in the block until this transaction was executed.
	CumulativeGasUsed *uint64 `json:"cumulative_gas_used,omitempty"`
	// The sum of the base fee and tip paid per unit of gas.
	EffectiveGasPrice *big.Int `json:"effective_gas_price,omitempty"`
	// Gas used by transaction
	GasUsed *uint64 `json:"gas_used,omitempty"`
	// Address of created contract if transaction was a contract creation
	ContractAddress *common.Address `json:"contract_address,omitempty"`
	// Bloom filter for logs produced by this transaction
	LogsBloom *BloomFilter `json:"logs_bloom,omitempty"`
	// Transaction type. For ethereum: Legacy, Eip2930, Eip1559, Eip4844
	Kind *uint8 `json:"type,omitempty"`
	// The Keccak 256-bit hash of the root node of the trie structure populated with each transaction in the transactions list portion of the block; formally Ht.
	Root *common.Hash `json:"root,omitempty"`
	// If transaction is executed successfully. This is the `statusCode`
	Status *uint8 `json:"status,omitempty"`
	// The fee associated with a transaction on the Layer 1, it is calculated as l1GasPrice multiplied by l1GasUsed
	L1Fee *big.Int `json:"l1_fee,omitempty"`
	// The gas price for transactions on the Layer 1
	L1GasPrice *big.Int `json:"l1_gas_price,omitempty"`
	// The amount of gas consumed by a transaction on the Layer 1
	L1GasUsed *uint64 `json:"l1_gas_used,omitempty"`
	// A multiplier applied to the actual gas usage on Layer 1 to calculate the dynamic costs. If set to 1, it has no impact on the L1 gas usage
	L1FeeScalar *float64 `json:"l1_fee_scalar,omitempty"`
	// Amount of gas spent on L1 calldata in units of L2 gas.
	GasUsedForL1 *uint64 `json:"gas_used_for_l1,omitempty"`
}

func NewTransactionFromRecord(schema *arrow.Schema, record arrow.Record) (*Transaction, error) {
	if record.NumCols() != int64(len(schema.Fields())) {
		return nil, errors.New("number of columns in record does not match schema")
	}

	toReturn := &Transaction{}

	for i, field := range schema.Fields() {
		col := record.Column(i)
		if col.Len() == 0 {
			continue
		}
		switch field.Name {
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
		case "from":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				from := common.BytesToAddress(val)
				toReturn.From = &from
			}
		case "gas":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.Gas = &val
			}
		case "gas_price":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.GasPrice = big.NewInt(0).SetUint64(val)
			}
		case "sighash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				hash := common.BytesToHash(val)
				toReturn.SigHash = &hash
			}
		case "hash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				hash := common.BytesToHash(val)
				toReturn.Hash = &hash
			}
		case "input":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				toReturn.Input = &val
			}
		case "nonce":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.Nonce = &val
			}
		case "to":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				to := common.BytesToAddress(val)
				toReturn.To = &to
			}
		case "transaction_index":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.TransactionIndex = &val
			}
		case "value":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.Value = big.NewInt(0).SetUint64(val)
			}
		case "v":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.V = big.NewInt(0).SetUint64(val)
			}
		case "r":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.R = big.NewInt(0).SetUint64(val)
			}
		case "s":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.S = big.NewInt(0).SetUint64(val)
			}
		case "y_parity":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.YParity = big.NewInt(0).SetUint64(val)
			}
		case "max_priority_fee_per_gas":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.MaxPriorityFeePerGas = big.NewInt(0).SetUint64(val)
			}
		case "max_fee_per_gas":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.MaxFeePerGas = big.NewInt(0).SetUint64(val)
			}
		case "chain_id":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.ChainID = big.NewInt(0).SetUint64(val)
			}
		case "access_list":
			if fCol, ok := col.(*array.List); ok {
				_ = fCol
			}
		case "max_fee_per_blob_gas":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.MaxFeePerBlobGas = big.NewInt(0).SetUint64(val)
			}
		case "blob_versioned_hashes":
			if fCol, ok := col.(*array.List); ok {
				_ = fCol
			}
		case "cumulative_gas_used":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.CumulativeGasUsed = &val
			}
		case "effective_gas_price":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.EffectiveGasPrice = big.NewInt(0).SetUint64(val)
			}
		case "gas_used":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.GasUsed = &val
			}
		case "contract_address":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				contractAddress := common.BytesToAddress(val)
				toReturn.ContractAddress = &contractAddress
			}
		case "logs_bloom":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				logsBloom := BloomFilter(val)
				toReturn.LogsBloom = &logsBloom
			}
		case "type":
			if fCol, ok := col.(*array.Uint8); ok {
				val := fCol.Value(0)
				toReturn.Kind = &val
			}
		case "root":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				root := common.BytesToHash(val)
				toReturn.Root = &root
			}
		case "status":
			if fCol, ok := col.(*array.Uint8); ok {
				val := fCol.Value(0)
				toReturn.Status = &val
			}
		case "l1_fee":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.L1Fee = big.NewInt(0).SetUint64(val)
			}
		case "l1_gas_price":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.L1GasPrice = big.NewInt(0).SetUint64(val)
			}
		case "l1_gas_used":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.L1GasUsed = &val
			}
		case "l1_fee_scalar":
			if fCol, ok := col.(*array.Float64); ok {
				val := fCol.Value(0)
				toReturn.L1FeeScalar = &val
			}
		case "gas_used_for_l1":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.GasUsedForL1 = &val
			}
		default:
			return nil, errors.New("unsupported field: " + field.Name)
		}
	}

	return toReturn, nil
}
