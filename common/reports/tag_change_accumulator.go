package reports

import (
	"bridgecrewio/yor/common/structure"
)

type TagChangeAccumulator struct {
	scannedBlocks      []structure.IBlock
	newTagTraces       []*ResourceRecord
	updatedTagTraces   []*ResourceRecord
	newBlockTraces     []structure.IBlock
	updatedBlockTraces []structure.IBlock
}

type ResourceRecord struct {
	file          string
	resource      string
	tagKey        string
	previousValue string
	newValue      string
	traceId       string
}

var TagChangeAccumulatorInstance *TagChangeAccumulator

func init() {
	TagChangeAccumulatorInstance = &TagChangeAccumulator{
		newTagTraces:     []*ResourceRecord{},
		updatedTagTraces: []*ResourceRecord{},
	}
}

func (a *TagChangeAccumulator) AccumulateChanges(block structure.IBlock) {
	a.scannedBlocks = append(a.scannedBlocks, block)
	diff := block.CalculateTagsDiff()
	// If only tags are new, add to newly traced. If some updates - add to updated. Otherwise will be added to
	// scannedBlocks.
	if len(diff.Updated) == 0 && len(diff.Added) > 0 {
		a.newBlockTraces = append(a.newBlockTraces, block)
	} else {
		if len(diff.Updated) > 0 {
			a.updatedBlockTraces = append(a.updatedBlockTraces, block)
		}
	}
	for _, tagDiff := range diff.Added {
		a.newTagTraces = append(a.newTagTraces, &ResourceRecord{
			file:          block.GetFilePath(),
			resource:      block.GetResourceId(),
			tagKey:        tagDiff.GetKey(),
			previousValue: "",
			newValue:      tagDiff.GetValue(),
			traceId:       block.GetTraceId(),
		})
	}

	for _, tagDiff := range diff.Updated {
		a.updatedTagTraces = append(a.updatedTagTraces, &ResourceRecord{
			file:          block.GetFilePath(),
			resource:      block.GetResourceId(),
			tagKey:        tagDiff.Key,
			previousValue: tagDiff.PrevValue,
			newValue:      tagDiff.NewValue,
			traceId:       block.GetTraceId(),
		})
	}
}

func (a *TagChangeAccumulator) GetTagChanges() ([]*ResourceRecord, []*ResourceRecord) {
	return a.newTagTraces, a.updatedTagTraces
}
func (a *TagChangeAccumulator) GetBlockChanges() ([]structure.IBlock, []structure.IBlock) {
	return a.newBlockTraces, a.updatedBlockTraces
}
func (a *TagChangeAccumulator) GetScannedBlocks() []structure.IBlock {
	return a.scannedBlocks
}
