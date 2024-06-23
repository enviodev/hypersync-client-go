package types

import (
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"math/big"
	"time"
)

// Block represents an Ethereum block object.
type Block struct {
	// A scalar value equal to the number of ancestor blocks. The genesis block has a number of zero; formally Hi.
	Number *big.Int `json:"number,omitempty"`
	// The Keccak 256-bit hash of the block
	Hash *common.Hash `json:"hash,omitempty"`
	// The Keccak 256-bit hash of the parent block’s header, in its entirety; formally Hp.
	ParentHash *common.Hash `json:"parent_hash,omitempty"`
	// A 64-bit value which, combined with the mixhash, proves that a sufficient amount of computation has been carried out on this block; formally Hn.
	Nonce *types.BlockNonce `json:"nonce,omitempty"`
	// The Keccak 256-bit hash of the ommers/uncles list portion of this block; formally Ho.
	Sha3Uncles *common.Hash `json:"sha3_uncles,omitempty"`
	// The Bloom filter composed from indexable information (logger address and log topics) contained in each log entry from the receipt of each transaction in the transactions list; formally Hb.
	LogsBloom *types.Bloom `json:"logs_bloom,omitempty"`
	// The Keccak 256-bit hash of the root node of the trie structure populated with each transaction in the transactions list portion of the block; formally Ht.
	TransactionsRoot *common.Hash `json:"transactions_root,omitempty"`
	// The Keccak 256-bit hash of the root node of the state trie, after all transactions are executed and finalisations applied; formally Hr.
	StateRoot *common.Hash `json:"state_root,omitempty"`
	// The Keccak 256-bit hash of the root node of the trie structure populated with each transaction in the transactions list portion of the block; formally Ht.
	ReceiptsRoot *common.Hash `json:"receipts_root,omitempty"`
	// The 160-bit address to which all fees collected from the successful mining of this block be transferred; formally Hc.
	Miner *common.Address `json:"miner,omitempty"`
	// A scalar value corresponding to the difficulty level of this block. This can be calculated from the previous block’s difficulty level and the timestamp; formally Hd.
	Difficulty *big.Int `json:"difficulty,omitempty"`
	// The cumulative sum of the difficulty of all blocks that have been mined in the Ethereum network since the inception of the network. It measures the overall security and integrity of the Ethereum network.
	TotalDifficulty *big.Int `json:"total_difficulty,omitempty"`
	// An arbitrary byte array containing data relevant to this block. This must be 32 bytes or fewer; formally Hx.
	ExtraData *[]byte `json:"extra_data,omitempty"`
	// The size of this block in bytes as an integer value, encoded as hexadecimal.
	Size *Quantity `json:"size,omitempty"`
	// A scalar value equal to the current limit of gas expenditure per block; formally Hl.
	GasLimit *uint64 `json:"gas_limit,omitempty"`
	// A scalar value equal to the total gas used in transactions in this block; formally Hg.
	GasUsed *uint64 `json:"gas_used,omitempty"`
	// A scalar value equal to the reasonable output of Unix’s time() at this block’s inception; formally Hs.
	Timestamp *time.Time `json:"timestamp,omitempty"`
	// Ommers/uncles header.
	Uncles *[]common.Hash `json:"uncles,omitempty"`
	// A scalar representing EIP1559 base fee which can move up or down each block according to a formula which is a function of gas used in parent block and gas target (block gas limit divided by elasticity multiplier) of parent block. The algorithm results in the base fee per gas increasing when blocks are above the gas target, and decreasing when blocks are below the gas target. The base fee per gas is burned.
	BaseFeePerGas *Quantity `json:"base_fee_per_gas,omitempty"`
	// The total amount of blob gas consumed by the transactions within the block, added in EIP-4844.
	BlobGasUsed *uint64 `json:"blob_gas_used,omitempty"`
	// A running total of blob gas consumed in excess of the target, prior to the block. Blocks with above-target blob gas consumption increase this value, blocks with below-target blob gas consumption decrease it (bounded at 0). This was added in EIP-4844.
	ExcessBlobGas *Quantity `json:"excess_blob_gas,omitempty"`
	// The hash of the parent beacon block's root is included in execution blocks, as proposed by EIP-4788. This enables trust-minimized access to consensus state, supporting staking pools, bridges, and more. The beacon roots contract handles root storage, enhancing Ethereum's functionalities.
	ParentBeaconBlockRoot *common.Hash `json:"parent_beacon_block_root,omitempty"`
	// The Keccak 256-bit hash of the withdrawals list portion of this block. See EIP-4895.
	WithdrawalsRoot *common.Hash `json:"withdrawals_root,omitempty"`
	// Withdrawal represents a validator withdrawal from the consensus layer.
	Withdrawals *[]Withdrawal `json:"withdrawals,omitempty"`
	// The L1 block number that would be used for block.number calls.
	L1BlockNumber *big.Int `json:"l1_block_number,omitempty"`
	// The number of L2 to L1 messages since Nitro genesis.
	SendCount *Quantity `json:"send_count,omitempty"`
	// The Merkle root of the outbox tree state.
	SendRoot *common.Hash `json:"send_root,omitempty"`
	// A 256-bit hash which, combined with the nonce, proves that a sufficient amount of computation has been carried out on this block; formally Hm.
	MixHash *common.Hash `json:"mix_hash,omitempty"`
}

func (b *Block) ToCommon() *types.Block {
	return &types.Block{}
}

func (b *Block) ToCommonHeader() *types.Header {
	return &types.Header{
		Number: b.Number,
		//Bloom:  *b.LogsBloom,
	}
}

func NewBlockFromRecord(schema *arrow.Schema, record arrow.Record) (*Block, error) {
	if record.NumCols() != int64(len(schema.Fields())) {
		return nil, errors.New("number of columns in record does not match schema")
	}

	toReturn := &Block{}

	for i, field := range schema.Fields() {
		col := record.Column(i)
		if col.Len() == 0 {
			continue
		}
		switch field.Name {
		case "number":
			if fCol, ok := col.(*array.Uint64); ok {
				val := fCol.Value(0)
				toReturn.Number = big.NewInt(0).SetUint64(val)
			}
		case "hash":
			if fCol, ok := col.(*array.Binary); ok {
				val := fCol.Value(0)
				hash := common.BytesToHash(val)
				toReturn.Hash = &hash
			}
		/*
			case "parent_hash":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.ParentHash = &Hash{val}
				}
			case "nonce":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.Nonce = &Nonce{val}
				}
			case "sha3_uncles":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.Sha3Uncles = &Hash{val}
				}
			case "logs_bloom":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.LogsBloom = &BloomFilter{val}
				}
			case "transactions_root":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.TransactionsRoot = &Hash{val}
				}
			case "state_root":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.StateRoot = &Hash{val}
				}
			case "receipts_root":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.ReceiptsRoot = &Hash{val}
				}
			case "miner":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.Miner = &Address{val}
				}
			case "difficulty":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.Difficulty = &Quantity{val}
				}
			case "total_difficulty":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.TotalDifficulty = &Quantity{val}
				}
			case "extra_data":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.ExtraData = &Data{val}
				}
			case "size":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.Size = &Quantity{val}
				}
			case "gas_limit":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.GasLimit = &Quantity{val}
				}
			case "gas_used":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.GasUsed = &Quantity{val}
				}
			case "timestamp":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.Timestamp = &Quantity{val}
				}
			case "uncles":
				if col, ok := col.(*array.List); ok {
					// Assuming uncles is a list of hashes
					uncles := make([]Hash, col.Len())
					for j := 0; j < col.Len(); j++ {
						uncle := col.Value(j).(*array.Binary).Value(0)
						uncles[j] = Hash{uncle}
					}
					toReturn.Uncles = &uncles
				}
			case "base_fee_per_gas":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.BaseFeePerGas = &Quantity{val}
				}
			case "blob_gas_used":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.BlobGasUsed = &Quantity{val}
				}
			case "excess_blob_gas":
				if col, ok := col.(*array.Uint64); ok {
					val := col.Value(0)
					toReturn.ExcessBlobGas = &Quantity{val}
				}
			case "parent_beacon_block_root":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.ParentBeaconBlockRoot = &Hash{val}
				}
			case "withdrawals_root":
				if col, ok := col.(*array.Binary); ok {
					val := col.Value(0)
					toReturn.WithdrawalsRoot = &Hash{val}
				}*/
		/*		case "withdrawals":
				if col, ok := col.(*array.List); ok {
					// Assuming withdrawals is a list of withdrawals
					withdrawals := make([]Withdrawal, col.Len())
					for j := 0; j < col.Len(); j++ {
						withdrawal := col.Value(j).(*array.Struct)
						// Populate withdrawal fields here
						withdrawals[j] = Withdrawal{}
					}
					toReturn.Withdrawals = &withdrawals
				}*/
		/*		case "l1_block_number":
					if col, ok := col.(*array.Uint64); ok {
						val := col.Value(0)
						toReturn.L1BlockNumber = &BlockNumber{val}
					}
				case "send_count":
					if col, ok := col.(*array.Uint64); ok {
						val := col.Value(0)
						toReturn.SendCount = &Quantity{val}
					}
				case "send_root":
					if col, ok := col.(*array.Binary); ok {
						val := col.Value(0)
						toReturn.SendRoot = &Hash{val}
					}
				case "mix_hash":
					if col, ok := col.(*array.Binary); ok {
						val := col.Value(0)
						toReturn.MixHash = &Hash{val}
					}*/
		default:
			return nil, errors.New("unsupported field: " + field.Name)
		}
	}

	return toReturn, nil
}
