package structure

import (
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"strings"
)

type TerraformBlock struct {
	structure.Block
	HclSyntaxBlock *hclsyntax.Block
}

func (b *TerraformBlock) GetResourceId() string {
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
func (b *TerraformBlock) GetLines() []int {
	r := b.HclSyntaxBlock.Body.Range()
	return []int{r.Start.Line, r.End.Line}
}

func (b *TerraformBlock) GetNewOwner() string {
	for _, tag := range b.GetNewTags() {
		if val, ok := tag.(*tags.GitModifiersTag); ok {
			return val.GetValue()
		}
	}
	return ""
}

func (b *TerraformBlock) GetPreviousOwner() string {
	for _, tag := range b.GetExistingTags() {
		if val, ok := tag.(*tags.GitModifiersTag); ok {
			return val.GetValue()
		}
	}
	return ""
}

func (b *TerraformBlock) GetTraceId() string {
	for _, tag := range b.GetExistingTags() {
		if val, ok := tag.(*tags.YorTraceTag); ok {
			return val.GetValue()
		}
	}
	for _, tag := range b.GetNewTags() {
		if val, ok := tag.(*tags.YorTraceTag); ok {
			return val.GetValue()
		}
	}
	return ""
}
