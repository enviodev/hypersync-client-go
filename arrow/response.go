package arrowhs

import "github.com/enviodev/hypersync-client-go/types"

type QueryResponseInterface interface {
	AppendBlockData(block types.Block)
	AppendTransactionData(block types.Transaction)
	GetData() types.DataResponse
	SetArchiveHeight(*int64)
	SetNextBlock(uint64)
	SetTotalExecutionTime(uint64)
	SetRollbackGuard(*types.RollbackGuard)
}
