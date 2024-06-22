package client

import (
	"bytes"
	"capnproto.org/go/capnp/v3"
	"fmt"
	"github.com/apache/arrow/go/v10/arrow/ipc"
	"github.com/davecgh/go-spew/spew"
	hypersynccapnp "github.com/enviodev/hypersync-client-go/capnp"
	"github.com/enviodev/hypersync-client-go/types"
	"github.com/enviodev/hypersync-client-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"io"
	"math/big"
	"os"
)

// Structs to represent the JSON structure
type Response struct {
	Data               []Data        `json:"data"`
	ArchiveHeight      int           `json:"archive_height"`
	NextBlock          int           `json:"next_block"`
	TotalExecutionTime int           `json:"total_execution_time"`
	RollbackGuard      RollbackGuard `json:"rollback_guard"`
}

type RollbackGuard struct {
	BlockNumber      int    `json:"block_number"`
	Timestamp        int64  `json:"timestamp"`
	Hash             string `json:"hash"`
	FirstBlockNumber int    `json:"first_block_number"`
	FirstParentHash  string `json:"first_parent_hash"`
}

type Transaction struct {
	BlockNumber *big.Int    `json:"block_number"`
	Hash        common.Hash `json:"hash"`
}

type Data struct {
	Transactions []*Transaction `json:"transactions"`
}

type ArrowResponse struct {
	ArchiveHeight      *int64 `json:"archive_height"`
	NextBlock          int64
	TotalExecutionTime int64
	//Data               ArrowResponseData
	RollbackGuard *RollbackGuard
}

type ArrowBatch struct{}

func parseQueryResponse(bReader io.ReadCloser) (*types.QueryResponse[*ArrowResponse], error) {
	defer bReader.Close()

	decoder := capnp.NewPackedDecoder(bReader)

	msg, err := decoder.Decode()
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode packed message")
	}

	queryResponse, err := hypersynccapnp.ReadRootQueryResponse(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get root pointer")
	}

	var archiveHeight *int64
	if queryResponse.ArchiveHeight() != -1 {
		h := queryResponse.ArchiveHeight()
		archiveHeight = &h
	}

	var rollbackGuard *types.RollbackGuard
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

		rollbackGuard = &types.RollbackGuard{
			BlockNumber:      rg.BlockNumber(),
			Timestamp:        rg.Timestamp(),
			Hash:             common.BytesToHash(hash),
			FirstBlockNumber: rg.FirstBlockNumber(),
			FirstParentHash:  common.BytesToHash(firstParentHash),
		}
	}

	dataPtr, dpErr := queryResponse.Data()
	if dpErr != nil {
		return nil, errors.Wrap(dpErr, "failed to read query response data")
	}

	if dataPtr.HasBlocks() {
		blocks, bErr := dataPtr.Blocks()
		if bErr != nil {
			return nil, errors.Wrap(bErr, "failed to parse block data")
		}

		blockData, bdErr := readChunks(blocks)
		if bdErr != nil {
			return nil, errors.Wrap(bdErr, "failed to read chunks from blocks data")
		}

		utils.DumpNodeWithExit(blockData)
	}

	toReturn := &types.QueryResponse[*ArrowResponse]{
		ArchiveHeight:      archiveHeight,
		NextBlock:          queryResponse.NextBlock(),
		TotalExecutionTime: queryResponse.TotalExecutionTime(),
		RollbackGuard:      rollbackGuard,
	}

	spew.Dump(toReturn)
	os.Exit(1)

	return nil, nil
}

func readChunks(data []byte) ([]ArrowBatch, error) {
	fmt.Println("Data Length:", len(data))
	fmt.Println("Data (hex):", fmt.Sprintf("%x", data))

	reader := bytes.NewReader(data)
	arrowReader, err := ipc.NewReader(reader, ipc.WithSchema(types.BlockHeaderSchema()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Arrow file reader")
	}

	schema := arrowReader.Schema()
	spew.Dump(schema)

	return nil, nil
}

/*func readChunks(data []byte) ([]ArrowBatch, error) {
	bufReader := memory.NewBufferBytes(data)
	arrowReader, err := ipc.NewFileReader(bufReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Arrow file reader")
	}
	defer arrowReader.Close()

	schema := arrowReader.Schema()
	var batches []ArrowBatch

	for i := 0; i < arrowReader.NumRecords(); i++ {
		record, err := arrowReader.ReadRecord(i)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read record batch")
		}
		batches = append(batches, ArrowBatch{
			Chunk:  record,
			Schema: schema,
		})
	}

	return batches, nil
}*/
