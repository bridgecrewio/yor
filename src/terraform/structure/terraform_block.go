package structure

import (
	"strings"

	"github.com/bridgecrewio/yor/src/common/structure"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type TerraformBlock struct {
	structure.Block
	HclSyntaxBlock *hclsyntax.Block
}

const ResourceBlockType = "resource"
const ModuleBlockType = "module"
const DataBlockType = "data"
const VarBlockType = "variable"

var SupportedBlockTypes = []string{ResourceBlockType, ModuleBlockType, VarBlockType}

func (b *TerraformBlock) GetResourceID() string {
	return strings.Join(b.HclSyntaxBlock.Labels, ".")
}

func (b *TerraformBlock) AddHclSyntaxBlock(hclSyntaxBlock *hclsyntax.Block) {
	b.HclSyntaxBlock = hclSyntaxBlock
}

func (b *TerraformBlock) GetLines(getContentLinesOnly ...bool) structure.Lines {
	r := b.HclSyntaxBlock.Body.Range()
	if len(getContentLinesOnly) == 0 || !getContentLinesOnly[0] {
		return structure.Lines{Start: r.Start.Line, End: r.End.Line}
	}

	endOfLastAttribute := r.Start.Line
	for _, attr := range b.HclSyntaxBlock.Body.Attributes {
		if attr.Range().End.Line > endOfLastAttribute {
			endOfLastAttribute = attr.Range().End.Line
		}
	}

	return structure.Lines{Start: r.Start.Line, End: endOfLastAttribute}
}

func (b *TerraformBlock) GetTagsLines() structure.Lines {
	for _, attr := range b.HclSyntaxBlock.Body.Attributes {
		if attr.Name == b.TagsAttributeName {
			return structure.Lines{Start: attr.SrcRange.Start.Line, End: attr.SrcRange.End.Line}
		}
	}
	return structure.Lines{Start: -1, End: -1}
}
func (b *TerraformBlock) GetSeparator() string {
	return "="
}
