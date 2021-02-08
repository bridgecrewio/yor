package reports

import (
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
)

type TagChangeAccumulator struct {
	changesByFile map[string]*structure.IBlock
	changesByTag  map[tags.Tag]*structure.IBlock
}

type ResourceRecord struct {
	file          string
	resource      string
	previousOwner string
	newOwner      string
	traceId       string
}

var accumulatorInstance *TagChangeAccumulator

func GetAccumulator() *TagChangeAccumulator {
	// get instance of singleton accumulator
	if accumulatorInstance == nil {
		accumulatorInstance = &TagChangeAccumulator{
			changesByFile: make(map[string]*structure.IBlock),
			changesByTag:  make(map[tags.Tag]*structure.IBlock),
		}
	}

	return accumulatorInstance
}

func (a *TagChangeAccumulator) AccumulateChanges(block *structure.IBlock) {
	// TODO
}

func (a *TagChangeAccumulator) GetPreviouslyTaggedResources() []*ResourceRecord {
	return nil
}
func (a *TagChangeAccumulator) GetUntaggedResources() []*ResourceRecord {
	return nil
}
func (a *TagChangeAccumulator) GetNewlyTaggedResources() []*ResourceRecord {
	return nil
}
