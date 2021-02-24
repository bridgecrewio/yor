package structure

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/structure"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type TerraformBlock struct {
	structure.Block
	HclSyntaxBlock *hclsyntax.Block
}

func (b *TerraformBlock) GetResourceID() string {
	return strings.Join(b.HclSyntaxBlock.Labels, ".")
}

func (b *TerraformBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *TerraformBlock) AddHclSyntaxBlock(hclSyntaxBlock *hclsyntax.Block) {
	b.HclSyntaxBlock = hclSyntaxBlock
}

func (b *TerraformBlock) String() string {
	// TODO
	return ""
}
func (b *TerraformBlock) GetLines() common.Lines {
	r := b.HclSyntaxBlock.Body.Range()
	return common.Lines{Start: r.Start.Line, End: r.End.Line}
}
