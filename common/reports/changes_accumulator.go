package reports

import (
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
)

type ChangesAccumulator struct {
	changesByFile map[string]*structure.Block
	changesByTag  map[tagging.Tag]*structure.Block
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

func (a *ChangesAccumulator) GetData() interface{} {
	// TODO - replace this method after the report structure is determined
	return nil
}
