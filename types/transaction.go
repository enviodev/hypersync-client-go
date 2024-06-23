package types

import (
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

type TransactionSelection struct {
	From            []Address `json:"from,omitempty"`
	To              []Address `json:"to,omitempty"`
	SigHash         []SigHash `json:"sighash,omitempty"`
	Status          *uint8    `json:"status,omitempty"`
	Kind            []uint8   `json:"type,omitempty"`
	ContractAddress []Address `json:"contract_address,omitempty"`
}

// Transaction represents an Ethereum transaction object.
type Transaction struct {
	// The Keccak 256-bit hash of the block
	BlockHash *common.Hash `json:"block_hash,omitempty"`
	// A scalar value equal to the number of ancestor blocks. The genesis block has a number of zero; formally Hi.
	BlockNumber *big.Int `json:"block_number,omitempty"`
	// The 160-bit address of the message call’s sender
	From *common.Address `json:"from,omitempty"`
	// A scalar value equal to the maximum amount of gas that should be used in executing this transaction. This is paid up-front, before any computation is done and may not be increased later; formally Tg.
	Gas *Quantity `json:"gas,omitempty"`
	// A scalar value equal to the number of Wei to be paid per unit of gas for all computation costs incurred as a result of the execution of this transaction; formally Tp.
	GasPrice *Quantity `json:"gas_price,omitempty"`
	// A transaction hash is a keccak hash of an RLP encoded signed transaction.
	Hash *common.Hash `json:"hash,omitempty"`
	// Input has two uses depending if transaction is Create or Call (if `to` field is None or Some). pub init: An unlimited size byte array specifying the EVM-code for the account initialisation procedure CREATE, data: An unlimited size byte array specifying the input data of the message call, formally Td.
	Input *Data `json:"input,omitempty"`
	// A scalar value equal to the number of transactions sent by the sender; formally Tn.
	Nonce *Quantity `json:"nonce,omitempty"`
	// The 160-bit address of the message call’s recipient or, for a contract creation transaction, ∅, used here to denote the only member of B0 ; formally Tt.
	To *Address `json:"to,omitempty"`
	// Index of the transaction in the block
	TransactionIndex *TransactionIndex `json:"transaction_index,omitempty"`
	// A scalar value equal to the number of Wei to be transferred to the message call’s recipient or, in the case of contract creation, as an endowment to the newly created account; formally Tv.
	Value *Quantity `json:"value,omitempty"`
	// Replay protection value based on chain_id. See EIP-155 for more info.
	V *Quantity `json:"v,omitempty"`
	// The R field of the signature; the point on the curve.
	R *Quantity `json:"r,omitempty"`
	// The S field of the signature; the point on the curve.
	S *Quantity `json:"s,omitempty"`
	// yParity: Signature Y parity; formally Ty
	YParity *Quantity `json:"y_parity,omitempty"`
	// Max Priority fee that transaction is paying As ethereum circulation is around 120mil eth as of 2022 that is around 120000000000000000000000000 wei we are safe to use u128 as its max number is: 340282366920938463463374607431768211455 This is also known as `GasTipCap`
	MaxPriorityFeePerGas *Quantity `json:"max_priority_fee_per_gas,omitempty"`
	// A scalar value equal to the maximum amount of gas that should be used in executing this transaction. This is paid up-front, before any computation is done and may not be increased later; formally Tg. As ethereum circulation is around 120mil eth as of 2022 that is around 120000000000000000000000000 wei we are safe to use u128 as its max number is: 340282366920938463463374607431768211455 This is also known as `GasFeeCap`
	MaxFeePerGas *Quantity `json:"max_fee_per_gas,omitempty"`
	// Added as EIP-155: Simple replay attack protection
	ChainID *Quantity `json:"chain_id,omitempty"`
	// The accessList specifies a list of addresses and storage keys; these addresses and storage keys are added into the `accessed_addresses` and `accessed_storage_keys` global sets (introduced in EIP-2929). A gas cost is charged, though at a discount relative to the cost of accessing outside the list.
	AccessList *[]AccessList `json:"access_list,omitempty"`
	// Max fee per data gas aka BlobFeeCap or blobGasFeeCap
	MaxFeePerBlobGas *Quantity `json:"max_fee_per_blob_gas,omitempty"`
	// It contains a vector of fixed size hash(32 bytes)
	BlobVersionedHashes *[]Hash `json:"blob_versioned_hashes,omitempty"`
	// The total amount of gas used in the block until this transaction was executed.
	CumulativeGasUsed *Quantity `json:"cumulative_gas_used,omitempty"`
	// The sum of the base fee and tip paid per unit of gas.
	EffectiveGasPrice *Quantity `json:"effective_gas_price,omitempty"`
	// Gas used by transaction
	GasUsed *Quantity `json:"gas_used,omitempty"`
	// Address of created contract if transaction was a contract creation
	ContractAddress *Address `json:"contract_address,omitempty"`
	// Bloom filter for logs produced by this transaction
	LogsBloom *BloomFilter `json:"logs_bloom,omitempty"`
	// Transaction type. For ethereum: Legacy, Eip2930, Eip1559, Eip4844
	Kind *TransactionType `json:"type,omitempty"`
	// The Keccak 256-bit hash of the root node of the trie structure populated with each transaction in the transactions list portion of the block; formally Ht.
	Root *Hash `json:"root,omitempty"`
	// If transaction is executed successfully. This is the `statusCode`
	Status *TransactionStatus `json:"status,omitempty"`
	// The fee associated with a transaction on the Layer 1, it is calculated as l1GasPrice multiplied by l1GasUsed
	L1Fee *Quantity `json:"l1_fee,omitempty"`
	// The gas price for transactions on the Layer 1
	L1GasPrice *Quantity `json:"l1_gas_price,omitempty"`
	// The amount of gas consumed by a transaction on the Layer 1
	L1GasUsed *Quantity `json:"l1_gas_used,omitempty"`
	// A multiplier applied to the actual gas usage on Layer 1 to calculate the dynamic costs. If set to 1, it has no impact on the L1 gas usage
	L1FeeScalar *float64 `json:"l1_fee_scalar,omitempty"`
	// Amount of gas spent on L1 calldata in units of L2 gas.
	GasUsedForL1 *Quantity `json:"gas_used_for_l1,omitempty"`
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
		case "hash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				hash := common.BytesToHash(val)
				toReturn.Hash = &hash
			}
		case "block_number":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.BlockNumber = big.NewInt(0).SetUint64(val)
			}

		default:
			return nil, errors.New("unsupported field: " + field.Name)
		}
	}

	return toReturn, nil
}
