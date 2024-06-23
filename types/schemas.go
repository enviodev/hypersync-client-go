package types

import "github.com/apache/arrow/go/v10/arrow"

func BlockHeaderSchema(metadata *arrow.Metadata) *arrow.Schema {
	fields := []arrow.Field{
		{Name: "number", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
		{Name: "hash", Type: hashDT(), Nullable: false},
		{Name: "parent_hash", Type: hashDT(), Nullable: false},
		{Name: "nonce", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "sha3_uncles", Type: hashDT(), Nullable: false},
		{Name: "logs_bloom", Type: arrow.BinaryTypes.Binary, Nullable: false},
		{Name: "transactions_root", Type: hashDT(), Nullable: false},
		{Name: "state_root", Type: hashDT(), Nullable: false},
		{Name: "receipts_root", Type: hashDT(), Nullable: false},
		{Name: "miner", Type: addrDT(), Nullable: false},
		{Name: "difficulty", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "total_difficulty", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "extra_data", Type: arrow.BinaryTypes.Binary, Nullable: false},
		{Name: "size", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "gas_limit", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "gas_used", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "timestamp", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "uncles", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "base_fee_per_gas", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "blob_gas_used", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "excess_blob_gas", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "parent_beacon_block_root", Type: hashDT(), Nullable: true},
		{Name: "withdrawals_root", Type: hashDT(), Nullable: true},
		{Name: "withdrawals", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "l1_block_number", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
		{Name: "send_count", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "send_root", Type: hashDT(), Nullable: true},
		{Name: "mix_hash", Type: hashDT(), Nullable: true},
	}
	return arrow.NewSchema(fields, metadata)
}

func BlockSchemaFieldsAsString() []string {
	toReturn := make([]string, 0)
	schema := BlockHeaderSchema(nil)
	for _, field := range schema.Fields() {
		toReturn = append(toReturn, field.Name)
	}
	return toReturn
}

func TransactionSchema() *arrow.Schema {
	fields := []arrow.Field{
		{Name: "block_hash", Type: hashDT(), Nullable: false},
		{Name: "block_number", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
		{Name: "from", Type: addrDT(), Nullable: true},
		{Name: "gas", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "gas_price", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "hash", Type: hashDT(), Nullable: false},
		{Name: "input", Type: arrow.BinaryTypes.Binary, Nullable: false},
		{Name: "nonce", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "to", Type: addrDT(), Nullable: true},
		{Name: "transaction_index", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
		{Name: "value", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "v", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "r", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "s", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "max_priority_fee_per_gas", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "max_fee_per_gas", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "chain_id", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "cumulative_gas_used", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "effective_gas_price", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "gas_used", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
		{Name: "contract_address", Type: addrDT(), Nullable: true},
		{Name: "logs_bloom", Type: arrow.BinaryTypes.Binary, Nullable: false},
		{Name: "type", Type: arrow.PrimitiveTypes.Uint8, Nullable: true},
		{Name: "root", Type: hashDT(), Nullable: true},
		{Name: "status", Type: arrow.PrimitiveTypes.Uint8, Nullable: true},
		{Name: "sighash", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "y_parity", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "access_list", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "l1_fee", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "l1_gas_price", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "l1_gas_used", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "l1_fee_scalar", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "gas_used_for_l1", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "max_fee_per_blob_gas", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "blob_versioned_hashes", Type: arrow.BinaryTypes.Binary, Nullable: true},
	}
	return arrow.NewSchema(fields, nil)
}

func LogSchema() *arrow.Schema {
	fields := []arrow.Field{
		{Name: "removed", Type: arrow.FixedWidthTypes.Boolean, Nullable: true},
		{Name: "log_index", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
		{Name: "transaction_index", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
		{Name: "transaction_hash", Type: hashDT(), Nullable: false},
		{Name: "block_hash", Type: hashDT(), Nullable: false},
		{Name: "block_number", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
		{Name: "address", Type: addrDT(), Nullable: false},
		{Name: "data", Type: arrow.BinaryTypes.Binary, Nullable: false},
		{Name: "topic0", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "topic1", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "topic2", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "topic3", Type: arrow.BinaryTypes.Binary, Nullable: true},
	}
	return arrow.NewSchema(fields, nil)
}

func TraceSchema() *arrow.Schema {
	fields := []arrow.Field{
		{Name: "from", Type: addrDT(), Nullable: true},
		{Name: "to", Type: addrDT(), Nullable: true},
		{Name: "call_type", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "gas", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "input", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "init", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "value", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "author", Type: addrDT(), Nullable: true},
		{Name: "reward_type", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "block_hash", Type: arrow.BinaryTypes.Binary, Nullable: false},
		{Name: "block_number", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
		{Name: "address", Type: addrDT(), Nullable: true},
		{Name: "code", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "gas_used", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "output", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "subtraces", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
		{Name: "trace_address", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "transaction_hash", Type: arrow.BinaryTypes.Binary, Nullable: true},
		{Name: "transaction_position", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
		{Name: "type", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "error", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "sighash", Type: arrow.BinaryTypes.Binary, Nullable: true},
	}
	return arrow.NewSchema(fields, nil)
}
