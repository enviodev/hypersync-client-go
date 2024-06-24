package streams

import (
	"math/big"
	"sync/atomic"
)

// BlockIterator manages the iteration over a range of blocks, supporting batching and
// thread-safe updates to the current offset.
type BlockIterator struct {
	offset    uint64  // current offset in the block range
	end       uint64  // end of the block range
	batchSize *uint64 // size of each batch of blocks
	stepGen   uint32  // generator for steps, currently unused
}

// NewBlockIterator creates a new BlockIterator with the specified offset, end, and batch size.
func NewBlockIterator(offset uint64, end uint64, step *uint64) *BlockIterator {
	return &BlockIterator{
		offset:    offset,
		end:       end,
		batchSize: step,
	}
}

// GetCurrentOffset returns the current offset in the block range.
func (b *BlockIterator) GetCurrentOffset() uint64 {
	return b.offset
}

// GetEnd returns the end of the block range.
func (b *BlockIterator) GetEnd() uint64 {
	return b.end
}

// GetEndAsBigInt returns the end of the block range as a big.Int.
func (b *BlockIterator) GetEndAsBigInt() *big.Int {
	return big.NewInt(0).SetUint64(b.end)
}

// Completed checks if the iteration has reached or passed the end of the block range.
func (b *BlockIterator) Completed() bool {
	return b.offset >= b.end
}

// Next returns the next batch of blocks to be processed. It updates the current offset and
// returns the start and end of the batch, as well as a boolean indicating if the operation
// was successful.
func (b *BlockIterator) Next() (start uint64, end uint64, ok bool) {
	if b.offset >= b.end {
		return 0, 0, false
	}

	start = b.offset
	step := atomic.LoadUint64(b.batchSize)
	b.offset = min(b.offset+step, b.end)
	return start, b.offset, true
}
