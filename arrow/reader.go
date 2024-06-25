package arrowhs

import (
	"bytes"
	"capnproto.org/go/capnp/v3"
	"fmt"
	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/ipc"
	hypersynccapnp "github.com/enviodev/hypersync-client-go/capnp"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"io"
	"math/big"
)

// Reader is responsible for reading and processing query responses.
type Reader struct {
	reader            io.Reader
	rootQueryResponse hypersynccapnp.QueryResponse
	response          QueryResponseInterface
}

// NewQueryResponseReader creates a new Reader instance from an io.ReadCloser.
func NewQueryResponseReader(bReader io.ReadCloser) (*Reader, error) {
	queryResponse := &types.QueryResponse{Data: types.DataResponse{}}
	return NewReader(bReader, queryResponse)
}

// NewReader initializes a Reader with a provided io.ReadCloser and QueryResponseInterface.
func NewReader(bReader io.ReadCloser, response QueryResponseInterface) (*Reader, error) {
	toReturn := &Reader{
		reader:   bReader,
		response: response,
	}

	decoder := capnp.NewPackedDecoder(bReader)
	msg, err := decoder.Decode()
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode packed message")
	}

	queryResponse, err := hypersynccapnp.ReadRootQueryResponse(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get root pointer")
	}
	toReturn.rootQueryResponse = queryResponse

	if queryResponse.ArchiveHeight() != -1 {
		toReturn.response.SetArchiveHeight(big.NewInt(queryResponse.ArchiveHeight()))
	}

	toReturn.response.SetNextBlock(big.NewInt(0).SetUint64(queryResponse.NextBlock()))
	toReturn.response.SetTotalExecutionTime(queryResponse.TotalExecutionTime())

	if queryResponse.HasRollbackGuard() {
		rg, rgErr := queryResponse.RollbackGuard()
		if rgErr != nil {
			return nil, errors.Wrap(rgErr, "failed to get rollback guard")
		}

		hash, hErr := rg.Hash()
		if hErr != nil {
			return nil, errors.Wrap(hErr, "failed to get rollback guard hash")
		}

		firstParentHash, fphErr := rg.FirstParentHash()
		if fphErr != nil {
			return nil, errors.Wrap(fphErr, "failed to get rollback guard first parent hash")
		}

		rollbackGuard := &types.RollbackGuard{
			BlockNumber:      big.NewInt(0).SetUint64(rg.BlockNumber()),
			Timestamp:        rg.Timestamp(),
			Hash:             common.BytesToHash(hash),
			FirstBlockNumber: rg.FirstBlockNumber(),
			FirstParentHash:  common.BytesToHash(firstParentHash),
		}
		toReturn.response.SetRollbackGuard(rollbackGuard)
	}

	if pdErr := toReturn.processData(); pdErr != nil {
		return nil, errors.Wrap(pdErr, "failed to process query response data")
	}

	return toReturn, nil
}

// GetRootQueryResponse returns the root query response.
func (r *Reader) GetRootQueryResponse() hypersynccapnp.QueryResponse {
	return r.rootQueryResponse
}

// GetQueryResponse returns the query response.
func (r *Reader) GetQueryResponse() *types.QueryResponse {
	return r.response.(*types.QueryResponse)
}

// processData processes the data in the query response.
func (r *Reader) processData() error {
	dataPtr, dpErr := r.rootQueryResponse.Data()
	if dpErr != nil {
		return errors.Wrap(dpErr, "failed to read query response data")
	}

	if dataPtr.HasBlocks() {
		blocks, bErr := dataPtr.Blocks()
		if bErr != nil {
			return errors.Wrap(bErr, "failed to parse block data")
		}

		if bdErr := r.readChunks(blocks, types.BlocksDataType); bdErr != nil {
			return errors.Wrap(bdErr, "failed to read chunks from blocks data")
		}
	}

	if dataPtr.HasTransactions() {
		blocks, bErr := dataPtr.Transactions()
		if bErr != nil {
			return errors.Wrap(bErr, "failed to parse transactions data")
		}

		if bdErr := r.readChunks(blocks, types.TransactionsDataType); bdErr != nil {
			return errors.Wrap(bdErr, "failed to read chunks from transactions data")
		}
	}

	if dataPtr.HasLogs() {
		blocks, bErr := dataPtr.Logs()
		if bErr != nil {
			return errors.Wrap(bErr, "failed to parse logs data")
		}

		if bdErr := r.readChunks(blocks, types.LogsDataType); bdErr != nil {
			return errors.Wrap(bdErr, "failed to read chunks from logs data")
		}
	}

	if dataPtr.HasTraces() {
		blocks, bErr := dataPtr.Traces()
		if bErr != nil {
			return errors.Wrap(bErr, "failed to parse traces data")
		}

		if bdErr := r.readChunks(blocks, types.TracesDataType); bdErr != nil {
			return errors.Wrap(bdErr, "failed to read chunks from traces data")
		}
	}

	return nil
}

// readChunks reads and processes chunks of data from the provided byte slice.
func (r *Reader) readChunks(data []byte, dt types.DataType) error {
	if len(data) < 16 { // Minimum length for a valid Arrow IPC message + Polaris Arrow
		return errors.New("data length is too short to be a valid Arrow IPC message")
	}

	// Strip the first 8 bytes (Polaris Arrow header)
	reader := bytes.NewBuffer(data[8:])
	arrowReader, err := ipc.NewReader(reader)
	if err != nil {
		return errors.Wrap(err, "failed to create arrow file reader")
	}

	rSchema := arrowReader.Schema()

	for arrowReader.Next() {
		bRec := arrowReader.Record()
		if bRec == nil {
			break
		}

		for i := int64(0); i < bRec.NumRows(); i++ {
			rec := bRec.NewSlice(i, i+1)
			if pbErr := r.processRecord(rec, rSchema, dt); pbErr != nil {
				return errors.Wrap(pbErr, "failed to process batch")
			}
		}
	}

	if arErr := arrowReader.Err(); arErr != nil {
		return errors.Wrap(arErr, "error encountered during reading")
	}

	return nil
}

// processRecord processes an Arrow record based on the provided data type.
func (r *Reader) processRecord(record arrow.Record, schema *arrow.Schema, dt types.DataType) error {
	switch dt {
	case types.BlocksDataType:
		if block, bErr := types.NewBlockFromRecord(schema, record); bErr != nil {
			return errors.Wrap(bErr, "failed to build block data from record")
		} else if block != nil {
			r.response.AppendBlockData(*block)
		}
	case types.TransactionsDataType:
		if tx, bErr := types.NewTransactionFromRecord(schema, record); bErr != nil {
			return errors.Wrap(bErr, "failed to build transaction data from record")
		} else if tx != nil {
			r.response.AppendTransactionData(*tx)
		}
	case types.LogsDataType:
		if log, bErr := types.NewLogFromRecord(schema, record); bErr != nil {
			return errors.Wrap(bErr, "failed to build log data from record")
		} else if log != nil {
			r.response.AppendLogData(*log)
		}
	case types.TracesDataType:
		if trace, bErr := types.NewTraceFromRecord(schema, record); bErr != nil {
			return errors.Wrap(bErr, "failed to build log data from record")
		} else if trace != nil {
			r.response.AppendTraceData(*trace)
		}
	default:
		return fmt.Errorf("unsupported data type %v", dt)
	}
	return nil
}

// Close closes the reader. Currently, it does nothing but is provided for future use.
func (r *Reader) Close() error {
	return nil
}
