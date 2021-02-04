package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	terraformStructure "bridgecrewio/yor/terraform/structure"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) IsBlockTaggable(block structure.IBlock) bool {
	terraformBlock, ok := block.(*terraformStructure.TerraformBlock)
	if !ok {
		return false
	}

	//TODO - implement + delete print
	print(terraformBlock)
	return true

}

func (t *TerraformTagger) CreateTagsForBlock(block structure.IBlock, gitBlame *git_service.GitBlame) {

}
