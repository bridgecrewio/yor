package structureutils

import (
	"bridgecrewio/yor/common"
	commonStructure "bridgecrewio/yor/common/structure"
)

type MockTestBlock struct {
	commonStructure.Block
}

func (b *MockTestBlock) Init(filePath string, rawBlock interface{}) {}

func (b *MockTestBlock) String() string {
	return ""
}

func (b *MockTestBlock) GetResourceID() string {
	return ""
}

func (b *MockTestBlock) GetLines() common.Lines {
	return common.Lines{Start: 1, End: 3}
}
