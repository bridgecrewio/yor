package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	"fmt"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) CreateTagsForBlock(block structure.IBlock, gitBlame *gitservice.GitBlame) {
	for _, tag := range t.Tags {
		err := tag.CalculateValue(gitBlame)
		if err != nil {
			logger.Logger.Warning(fmt.Sprintf("failed to calculate tag value of tag %v, err: %s", tag, err))
			continue
		}
	}
	block.AddNewTags(t.Tags)
}
