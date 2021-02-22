package structureutils

import commonStructure "bridgecrewio/yor/common/structure"

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

func (b *MockTestBlock) GetLines() []int {
	return []int{0, 2}
}
