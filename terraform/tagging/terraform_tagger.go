package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) CreateTagsForBlock(block structure.IBlock, gitBlame *git_service.GitBlame) {
	print(block, gitBlame)
}
