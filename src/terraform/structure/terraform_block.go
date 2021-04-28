package structure

import (
	"bridgecrewio/yor/src/common"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type TerraformBlock struct {
	common.Block
	HclSyntaxBlock *hclsyntax.Block
}

func (b *TerraformBlock) UpdateTags() {
	return
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
	if len(getContentLinesOnly) == 0 || !getContentLinesOnly[0] {
		return common.Lines{Start: r.Start.Line, End: r.End.Line}
	}

	endOfLastAttribute := r.Start.Line
	for _, attr := range b.HclSyntaxBlock.Body.Attributes {
		if attr.Range().End.Line > endOfLastAttribute {
			endOfLastAttribute = attr.Range().End.Line
		}
	}

	return common.Lines{Start: r.Start.Line, End: endOfLastAttribute}
}

func (b *TerraformBlock) GetTagsLines() common.Lines {
	for _, attr := range b.HclSyntaxBlock.Body.Attributes {
		if attr.Name == b.TagsAttributeName {
			return common.Lines{Start: attr.SrcRange.Start.Line, End: attr.SrcRange.End.Line}
		}
	}
	return common.Lines{Start: -1, End: -1}
}
func (b *TerraformBlock) GetSeparator() string {
	return "="
}
