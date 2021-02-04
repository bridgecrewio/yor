package reports

import (
	"bridgecrewio/yor/structure"
	"bridgecrewio/yor/tagging"
)

type Accumulator struct {
	changesByFile map[string]*structure.Block
	changesByTag  map[tagging.Tag]*structure.Block
}

var accumulatorInstance *Accumulator

func GetAccumulator() *Accumulator {
	// get instance of singleton accumulator
	if accumulatorInstance == nil {
		accumulatorInstance = &Accumulator{
			changesByFile: make(map[string]*structure.Block),
			changesByTag:  make(map[tagging.Tag]*structure.Block),
		}
	}

	return accumulatorInstance
}

func (a *Accumulator) AccumulateChanges(block *structure.Block) {
	// TODO
}

func (a *Accumulator) GetData() interface{} {
	// TODO - replace this method after the report structure is determined
	return nil
}
