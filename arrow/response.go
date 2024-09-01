package arrowhs

import (
	"math/big"

	"github.com/markovichecha/hypersync-client-go/types"
)

// QueryResponseInterface defines the interface for handling query response data.
type QueryResponseInterface interface {
	// AppendBlockData appends a block to the query response data.
	AppendBlockData(block types.Block)

	// AppendTransactionData appends a transaction to the query response data.
	AppendTransactionData(transaction types.Transaction)

	// AppendLogData appends a log to the query response data.
	AppendLogData(log types.Log)

	// AppendTraceData appends a trace to the query response data.
	AppendTraceData(trace types.Trace)

	// GetData returns the data response.
	GetData() types.DataResponse

	// SetArchiveHeight sets the archive height in the query response.
	SetArchiveHeight(*big.Int)

	// SetNextBlock sets the next block in the query response.
	SetNextBlock(*big.Int)

	// HasNextBlock checks if the next block is set in the query response.
	HasNextBlock() bool

	// SetTotalExecutionTime sets the total execution time in the query response.
	SetTotalExecutionTime(uint64)

	// HasRollbackGuard checks if the rollback guard is set in the query response.
	HasRollbackGuard() bool

	// SetRollbackGuard sets the rollback guard in the query response.
	SetRollbackGuard(*types.RollbackGuard)

	// GetRollbackGuard returns the rollback guard from the query response.
	GetRollbackGuard() *types.RollbackGuard
}
