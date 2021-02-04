package reports

import (
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
)

type ChangesAccumulator struct {
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

var accumulatorInstance *ChangesAccumulator

func GetAccumulator() *ChangesAccumulator {
	// get instance of singleton accumulator
	if accumulatorInstance == nil {
		accumulatorInstance = &ChangesAccumulator{
			changesByFile: make(map[string]*structure.Block),
			changesByTag:  make(map[tagging.Tag]*structure.Block),
		}
	}

	return accumulatorInstance
}

func (a *ChangesAccumulator) AccumulateChanges(block *structure.Block) {
	// TODO
}

func (a *ChangesAccumulator) GetPreviouslyTaggedResources() []*ResourceRecord {
	return nil
}
func (a *ChangesAccumulator) GetUntaggedResources() []*ResourceRecord {
	return nil
}
func (a *ChangesAccumulator) GetNewlyTaggedResources() []*ResourceRecord {
	return nil
}
