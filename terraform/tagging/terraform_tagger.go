package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/tagging"
	terraformStructure "bridgecrewio/yor/terraform/structure"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) IsBlockTaggable(block *terraformStructure.TerraformBlock) bool {
	//TODO - implement + delete print
	print(block)
	return true

}

func (t *TerraformTagger) CreateTagsForBlock(block *terraformStructure.TerraformBlock, gitBlame *git_service.GitBlame) {
	_, _ = block, gitBlame

}
