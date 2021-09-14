package reports

import (
	"sync"

	"github.com/bridgecrewio/yor/src/common/structure"
)

type TagChangeAccumulator struct {
	ScannedBlocks      []structure.IBlock
	NewBlockTraces     []structure.IBlock
	UpdatedBlockTraces []structure.IBlock
}

var TagChangeAccumulatorInstance *TagChangeAccumulator
var accumulatorLock sync.Mutex

func init() {
	TagChangeAccumulatorInstance = &TagChangeAccumulator{}
}

// AccumulateChanges saves the results of the scan of each block.
// If a block has no changes, it will be saved only to ScannedBlocks
// Otherwise it will be saved to NewBlockTraces if it is new or to UpdatedBlockTraces otherwise
func (a *TagChangeAccumulator) AccumulateChanges(block structure.IBlock) {
	accumulatorLock.Lock()
	defer accumulatorLock.Unlock()
	a.ScannedBlocks = append(a.ScannedBlocks, block)
	diff := block.CalculateTagsDiff()
	// If only tags are new, add to newly traced. If some updates - add to updated. Otherwise will be added to
	// ScannedBlocks.
	if len(diff.Updated) == 0 && len(diff.Added) > 0 {
		a.NewBlockTraces = append(a.NewBlockTraces, block)
	} else if len(diff.Updated) > 0 {
		a.UpdatedBlockTraces = append(a.UpdatedBlockTraces, block)
	}
}

// GetBlockChanges returns both the NewBlockTraces and the UpdatedBlockTraces that were found by the parsers
func (a *TagChangeAccumulator) GetBlockChanges() ([]structure.IBlock, []structure.IBlock) {
	return a.NewBlockTraces, a.UpdatedBlockTraces
}

func (a *TagChangeAccumulator) GetScannedBlocks() []structure.IBlock {
	return a.ScannedBlocks
}
