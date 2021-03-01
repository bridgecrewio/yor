package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/structure"
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

func (b *TerraformBlock) GetLines(getContentLinesOnly ...bool) common.Lines {
	r := b.HclSyntaxBlock.Body.Range()
	if len(getContentLinesOnly) == 0 || getContentLinesOnly[0] == false {
		return common.Lines{Start: r.Start.Line, End: r.End.Line}
	} else {
		endOfLastAttribute := r.Start.Line
		for _, attr := range b.HclSyntaxBlock.Body.Attributes {
			if attr.Range().End.Line > endOfLastAttribute {
				endOfLastAttribute = attr.Range().End.Line
			}
		}

		return common.Lines{Start: r.Start.Line, End: endOfLastAttribute}
	}
}
