package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
	"strings"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) CreateTagsForBlock(block structure.IBlock, gitBlame *gitservice.GitBlame) {
	var newTags []tags.ITag
	for _, tag := range t.Tags {
		err := tag.CalculateValue(gitBlame)
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to calculate tag value of tag %v, err: %s", tag, err))
			continue
		}
		newTags = append(newTags, &tags.Tag{Key: tag.GetKey(), Value: tag.GetValue()})
	}
	block.AddNewTags(newTags)
}

func (t *TerraformTagger) IsFileSkipped(file string) bool {
	ignoredPatterns := []string{".terraform"}
	ignoredDirs := t.GetSkippedDirs()
	ignoredPatterns = append(ignoredPatterns, ignoredDirs...)
	for _, pattern := range ignoredPatterns {
		if strings.Contains(file, pattern) {
			return true
		}
	}
	return false
}
