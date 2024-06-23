package options

import "math/big"

// HexOutput represents the formatting of binary columns numbers into UTF8 hex.
type HexOutput int

const (
	HexOutputDefault HexOutput = iota
)

// ColumnMapping represents the column mapping for stream function output.
// TODO
type ColumnMapping struct{}

// StreamOptions represents the configuration for hypersync event streaming.
type StreamOptions struct {
	// ColumnMapping is the column mapping for stream function output.
	// It lets you map columns you want into the DataTypes you want.
	ColumnMapping *ColumnMapping `mapstructure:"columnMapping" yaml:"columnMapping" json:"columnMapping"`

	// EventSignature is used to populate decode logs. Decode logs would be empty if set to None.
	EventSignature *string `mapstructure:"eventSignature" yaml:"eventSignature" json:"eventSignature"`

	// HexOutput determines the formatting of binary columns numbers into UTF8 hex.
	HexOutput HexOutput `mapstructure:"hexOutput" yaml:"hexOutput" json:"hexOutput"`

	// BatchSize is the initial batch size. Size would be adjusted based on response size during execution.
	BatchSize *uint64 `mapstructure:"batchSize" yaml:"batchSize" json:"batchSize"`

	// MaxBatchSize is the maximum batch size that could be used during dynamic adjustment.
	MaxBatchSize *big.Int `mapstructure:"maxBatchSize" yaml:"maxBatchSize" json:"maxBatchSize"`

	// MinBatchSize is the minimum batch size that could be used during dynamic adjustment.
	MinBatchSize *big.Int `mapstructure:"minBatchSize" yaml:"minBatchSize" json:"minBatchSize"`

	// Concurrency is the number of async threads that would be spawned to execute different block ranges of queries.
	Concurrency *big.Int `mapstructure:"concurrency" yaml:"concurrency" json:"concurrency"`

	// MaxNumBlocks is the max number of blocks to fetch in a single request.
	MaxNumBlocks *big.Int `mapstructure:"maxNumBlocks" yaml:"maxNumBlocks" json:"maxNumBlocks"`

	// MaxNumTransactions is the max number of transactions to fetch in a single request.
	MaxNumTransactions *big.Int `mapstructure:"maxNumTransactions" yaml:"maxNumTransactions" json:"maxNumTransactions"`

	// MaxNumLogs is the max number of logs to fetch in a single request.
	MaxNumLogs *big.Int `mapstructure:"maxNumLogs" yaml:"maxNumLogs" json:"maxNumLogs"`

	// MaxNumTraces is the max number of traces to fetch in a single request.
	MaxNumTraces *big.Int `mapstructure:"maxNumTraces" yaml:"maxNumTraces" json:"maxNumTraces"`

	// ResponseBytesCeiling is the size of a response in bytes from which step size will be lowered.
	ResponseBytesCeiling *big.Int `mapstructure:"responseBytesCeiling" yaml:"responseBytesCeiling" json:"responseBytesCeiling"`

	// ResponseBytesFloor is the size of a response in bytes from which step size will be increased.
	ResponseBytesFloor *big.Int `mapstructure:"responseBytesFloor" yaml:"responseBytesFloor" json:"responseBytesFloor"`
}

func DefaultStreamOptions() *StreamOptions {
	return &StreamOptions{
		Concurrency:  big.NewInt(10),
		MaxBatchSize: big.NewInt(1024),
	}
}
