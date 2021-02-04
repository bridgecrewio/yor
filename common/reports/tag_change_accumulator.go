package reports

import (
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
)

type TagChangeAccumulator struct {
	changesByFile map[string]*structure.Block
	changesByTag  map[tagging.Tag]*structure.Block
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
			changesByFile: make(map[string]*structure.Block),
			changesByTag:  make(map[tagging.Tag]*structure.Block),
		}
	}

	return accumulatorInstance
}

func (a *TagChangeAccumulator) AccumulateChanges(block *structure.Block) {
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
