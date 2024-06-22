package types

// Block represents an Ethereum block object.
type Block struct {
	// A scalar value equal to the number of ancestor blocks. The genesis block has a number of zero; formally Hi.
	Number *uint64 `json:"number,omitempty"`
	// The Keccak 256-bit hash of the block
	Hash *Hash `json:"hash,omitempty"`
	// The Keccak 256-bit hash of the parent block’s header, in its entirety; formally Hp.
	ParentHash *Hash `json:"parent_hash,omitempty"`
	// A 64-bit value which, combined with the mixhash, proves that a sufficient amount of computation has been carried out on this block; formally Hn.
	Nonce *Nonce `json:"nonce,omitempty"`
	// The Keccak 256-bit hash of the ommers/uncles list portion of this block; formally Ho.
	Sha3Uncles *Hash `json:"sha3_uncles,omitempty"`
	// The Bloom filter composed from indexable information (logger address and log topics) contained in each log entry from the receipt of each transaction in the transactions list; formally Hb.
	LogsBloom *BloomFilter `json:"logs_bloom,omitempty"`
	// The Keccak 256-bit hash of the root node of the trie structure populated with each transaction in the transactions list portion of the block; formally Ht.
	TransactionsRoot *Hash `json:"transactions_root,omitempty"`
	// The Keccak 256-bit hash of the root node of the state trie, after all transactions are executed and finalisations applied; formally Hr.
	StateRoot *Hash `json:"state_root,omitempty"`
	// The Keccak 256-bit hash of the root node of the trie structure populated with each transaction in the transactions list portion of the block; formally Ht.
	ReceiptsRoot *Hash `json:"receipts_root,omitempty"`
	// The 160-bit address to which all fees collected from the successful mining of this block be transferred; formally Hc.
	Miner *Address `json:"miner,omitempty"`
	// A scalar value corresponding to the difficulty level of this block. This can be calculated from the previous block’s difficulty level and the timestamp; formally Hd.
	Difficulty *Quantity `json:"difficulty,omitempty"`
	// The cumulative sum of the difficulty of all blocks that have been mined in the Ethereum network since the inception of the network. It measures the overall security and integrity of the Ethereum network.
	TotalDifficulty *Quantity `json:"total_difficulty,omitempty"`
	// An arbitrary byte array containing data relevant to this block. This must be 32 bytes or fewer; formally Hx.
	ExtraData *Data `json:"extra_data,omitempty"`
	// The size of this block in bytes as an integer value, encoded as hexadecimal.
	Size *Quantity `json:"size,omitempty"`
	// A scalar value equal to the current limit of gas expenditure per block; formally Hl.
	GasLimit *Quantity `json:"gas_limit,omitempty"`
	// A scalar value equal to the total gas used in transactions in this block; formally Hg.
	GasUsed *Quantity `json:"gas_used,omitempty"`
	// A scalar value equal to the reasonable output of Unix’s time() at this block’s inception; formally Hs.
	Timestamp *Quantity `json:"timestamp,omitempty"`
	// Ommers/uncles header.
	Uncles *[]Hash `json:"uncles,omitempty"`
	// A scalar representing EIP1559 base fee which can move up or down each block according to a formula which is a function of gas used in parent block and gas target (block gas limit divided by elasticity multiplier) of parent block. The algorithm results in the base fee per gas increasing when blocks are above the gas target, and decreasing when blocks are below the gas target. The base fee per gas is burned.
	BaseFeePerGas *Quantity `json:"base_fee_per_gas,omitempty"`
	// The total amount of blob gas consumed by the transactions within the block, added in EIP-4844.
	BlobGasUsed *Quantity `json:"blob_gas_used,omitempty"`
	// A running total of blob gas consumed in excess of the target, prior to the block. Blocks with above-target blob gas consumption increase this value, blocks with below-target blob gas consumption decrease it (bounded at 0). This was added in EIP-4844.
	ExcessBlobGas *Quantity `json:"excess_blob_gas,omitempty"`
	// The hash of the parent beacon block's root is included in execution blocks, as proposed by EIP-4788. This enables trust-minimized access to consensus state, supporting staking pools, bridges, and more. The beacon roots contract handles root storage, enhancing Ethereum's functionalities.
	ParentBeaconBlockRoot *Hash `json:"parent_beacon_block_root,omitempty"`
	// The Keccak 256-bit hash of the withdrawals list portion of this block. See EIP-4895.
	WithdrawalsRoot *Hash `json:"withdrawals_root,omitempty"`
	// Withdrawal represents a validator withdrawal from the consensus layer.
	Withdrawals *[]Withdrawal `json:"withdrawals,omitempty"`
	// The L1 block number that would be used for block.number calls.
	L1BlockNumber *BlockNumber `json:"l1_block_number,omitempty"`
	// The number of L2 to L1 messages since Nitro genesis.
	SendCount *Quantity `json:"send_count,omitempty"`
	// The Merkle root of the outbox tree state.
	SendRoot *Hash `json:"send_root,omitempty"`
	// A 256-bit hash which, combined with the nonce, proves that a sufficient amount of computation has been carried out on this block; formally Hm.
	MixHash *Hash `json:"mix_hash,omitempty"`
}
