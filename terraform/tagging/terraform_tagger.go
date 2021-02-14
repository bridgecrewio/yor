package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	"fmt"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) CreateTagsForBlock(block structure.IBlock, gitBlame *gitservice.GitBlame) error {
	for _, tag := range t.Tags {
		err := tag.CalculateValue(gitBlame)
		if err != nil {
			return fmt.Errorf("failed to calculate tag value of tag %v, err: %s", tag, err)
		}
	}
	block.AddNewTags(t.Tags)
	return nil
}
