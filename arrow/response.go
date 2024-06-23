package arrowhs

import (
	"github.com/enviodev/hypersync-client-go/types"
	"math/big"
)

type QueryResponseInterface interface {
	AppendBlockData(block types.Block)
	AppendTransactionData(block types.Transaction)
	GetData() types.DataResponse
	SetArchiveHeight(*big.Int)
	SetNextBlock(*big.Int)
	SetTotalExecutionTime(uint64)
	SetRollbackGuard(*types.RollbackGuard)
}
