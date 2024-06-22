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
)

type Reader struct {
	reader            io.Reader
	rootQueryResponse hypersynccapnp.QueryResponse
	response          QueryResponseInterface
}

func NewQueryResponseReader(bReader io.ReadCloser) (*Reader, error) {
	queryResponse := &types.QueryResponse{Data: types.DataResponse{}}
	return NewReader(bReader, queryResponse)
}

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

	var archiveHeight *int64
	if queryResponse.ArchiveHeight() != -1 {
		h := queryResponse.ArchiveHeight()
		archiveHeight = &h
	}
	toReturn.response.SetArchiveHeight(archiveHeight)
	toReturn.response.SetNextBlock(queryResponse.NextBlock())
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
			BlockNumber:      rg.BlockNumber(),
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

func (r *Reader) GetRootQueryResponse() hypersynccapnp.QueryResponse {
	return r.rootQueryResponse
}

func (r *Reader) GetQueryResponse() *types.QueryResponse {
	return r.response.(*types.QueryResponse)
}

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

	return nil
}

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
		rec := arrowReader.Record()
		if rec == nil {
			break
		}

		fmt.Println("readChunks - Record NumCols:", rec.NumCols())

		if pbErr := r.processRecord(rec, rSchema, dt); pbErr != nil {
			return errors.Wrap(pbErr, "failed to process batch")
		}

	}

	if arErr := arrowReader.Err(); arErr != nil {
		return errors.Wrap(arErr, "error encountered during reading")
	}

	return nil
}

func (r *Reader) processRecord(record arrow.Record, schema *arrow.Schema, dt types.DataType) error {
	fmt.Println("Process batch record columns:", record.NumCols())

	switch dt {
	case types.BlocksDataType:
		r.response.AppendBlockData(types.Block{})
	case types.TransactionsDataType:
		r.response.AppendTransactionData(types.Transaction{})
	default:
		return fmt.Errorf("unsupported data type %v", dt)
	}
	return nil
}

func (r *Reader) Close() error {
	return nil
}
