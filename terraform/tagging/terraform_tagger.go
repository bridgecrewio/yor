package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	tfStructure "bridgecrewio/yor/terraform/structure"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) IsBlockTaggable(block interface{}) bool {
	tfBlock, ok := block.(*tfStructure.TerraformBlock)
	if !ok {
		return false
	}
	//TODO - implement + delete print
	print(tfBlock)
	return true

}

func (t *TerraformTagger) CreateTagsForBlock(block structure.IBlock, gitBlame *git_service.GitBlame) {
	tfBlock, ok := block.(*tfStructure.TerraformBlock)
	if !ok {
		return
	}
	//TODO - implement + delete print
	print(tfBlock, gitBlame)
}
