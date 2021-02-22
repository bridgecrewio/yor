package structure

import (
	"bridgecrewio/yor/common/structure"
)

type CloudformationBlock struct {
	structure.Block
	lines []int
	name  string
}

func (b *CloudformationBlock) GetResourceID() string {
	return b.name
}

func (b *CloudformationBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *CloudformationBlock) String() string {
	// TODO
	return ""
}
func (b *CloudformationBlock) GetLines() []int {
	return b.lines
}
