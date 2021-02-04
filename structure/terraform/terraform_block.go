package terraform

import "bridgecrewio/yor/structure"

type TerraformBlock struct {
	structure.Block
}

func (b *TerraformBlock) Init(filePath string, rawBlock interface{}) {
	// TODO
}

func (b *TerraformBlock) String() string {
	// TODO
	return ""
}
func (b *TerraformBlock) GetLines() []int {
	// TODO
	return nil
}
