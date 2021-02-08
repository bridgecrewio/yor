package structure

import (
	"bridgecrewio/yor/common/structure"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type TerraformBlock struct {
	structure.Block
	hclSyntaxBlock *hclsyntax.Block
	NewOwner       string
	PreviousOwner  string
	TraceId        string
}

func (b *TerraformBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *TerraformBlock) AddHclSyntaxBlock(hclSyntaxBlock *hclsyntax.Block) {
	b.hclSyntaxBlock = hclSyntaxBlock
}

func (b *TerraformBlock) String() string {
	// TODO
	return ""
}
func (b *TerraformBlock) GetLines() []int {
	r := b.hclSyntaxBlock.Body.Range()
	return []int{r.Start.Line, r.End.Line}
}

func (b *TerraformBlock) GetNewOwner() string {
	return b.NewOwner
}

func (b *TerraformBlock) GetPreviousOwner() string {
	return b.PreviousOwner
}

func (b *TerraformBlock) GetTraceId() string {
	return b.TraceId
}
