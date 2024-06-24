package options

import (
	"fmt"
	"math/big"
	"runtime"
)

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

	// HexOutput determines the formatting of binary columns numbers into UTF8 hex.
	HexOutput HexOutput `mapstructure:"hexOutput" yaml:"hexOutput" json:"hexOutput"`

	// Concurrency is the number of async threads that would be spawned to execute different block ranges of queries.
	Concurrency *big.Int `mapstructure:"concurrency" yaml:"concurrency" json:"concurrency"`

	// BatchSize is the initial batch size. Size would be adjusted based on response size during execution.
	BatchSize *big.Int `mapstructure:"batchSize" yaml:"batchSize" json:"batchSize"`

	// MaxBatchSize is the maximum batch size that could be used during dynamic adjustment.
	MaxBatchSize *big.Int `mapstructure:"maxBatchSize" yaml:"maxBatchSize" json:"maxBatchSize"`

	// MinBatchSize is the minimum batch size that could be used during dynamic adjustment.
	MinBatchSize *big.Int `mapstructure:"minBatchSize" yaml:"minBatchSize" json:"minBatchSize"`
}

func (s *StreamOptions) Validate() error {
	if s.Concurrency == nil || s.Concurrency.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid stream concurrency provided")
	}
	if s.BatchSize == nil || s.BatchSize.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid stream batch size provided")
	}
	return nil
}

func DefaultStreamOptions() *StreamOptions {
	return &StreamOptions{
		Concurrency: big.NewInt(0).SetInt64(int64(runtime.NumCPU())),
		BatchSize:   big.NewInt(4096),
	}
}

func DefaultStreamOptionsWithBatchSize(batchSize *big.Int) *StreamOptions {
	return &StreamOptions{
		Concurrency: big.NewInt(0).SetInt64(int64(runtime.NumCPU())),
		BatchSize:   batchSize,
	}
}
