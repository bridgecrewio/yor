package structure

import (
	"bridgecrewio/yor/common/structure"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type TerraformBlock struct {
	structure.Block
}

func (b *TerraformBlock) Init(filePath string, rawBlock interface{}) error {
	b.RawBlock = rawBlock
	b.FilePath = filePath
	//b.GetLines()

	return nil
}

func (b *TerraformBlock) String() string {
	// TODO
	return ""
}
func (b *TerraformBlock) GetLines() []int {
	hclBlock := b.RawBlock.(*hclwrite.Block)
	//hclBlock.BuildTokens()
	print(hclBlock)

	return nil
}
