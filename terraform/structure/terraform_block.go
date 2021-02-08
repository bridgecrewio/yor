package structure

import "bridgecrewio/yor/common/structure"

type TerraformBlock struct {
	structure.Block
}

func (b *TerraformBlock) Init(filePath string, rawBlock interface{}) {
	b.Block.FilePath = filePath
	b.RawBlock = rawBlock
}

func (b *TerraformBlock) String() string {
	// TODO
	return ""
}
func (b *TerraformBlock) GetLines() []int {
	// TODO
	return nil
}

func (b *TerraformBlock) GetRawBlock() interface{} {
	return nil
}

func (b *TerraformBlock) GetNewOwner() string {
	return ""
}

func (b *TerraformBlock) GetPreviousOwner() string {
	return ""
}

func (b *TerraformBlock) GetTraceId() string {
	panic("implement me")
}

func (b *TerraformBlock) IsTaggable() bool {
	panic("implement me")
}
