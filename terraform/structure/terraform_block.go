package structure

import (
	"bridgecrewio/yor/common/structure"
)

type TerraformBlock struct {
	structure.Block
}

func (b *TerraformBlock) Init(filePath string, rawBlock interface{}) error {
	b.RawBlock = rawBlock
	b.FilePath = filePath

	return nil
}

func (b *TerraformBlock) String() string {
	// TODO
	return ""
}
func (b *TerraformBlock) GetLines() []int {
	// TODO
	return nil
}
