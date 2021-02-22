package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
)

type GitTagger struct {
	Tagger
	gitService *gitservice.GitService
}

func (t *GitTagger) InitTagger(path string) {
	gitService, err := gitservice.NewGitService(path)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize git service for path %s", path))
	}
	t.gitService = gitService
}

func (t *GitTagger) CreateTagsForBlock(block structure.IBlock) {
	blame, err := t.gitService.GetBlameForFileLines(block.GetFilePath(), block.GetLines())
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to tag %v with git tags, err: %v", block.GetResourceID(), err.Error()))
		return
	}
	if blame == nil {
		logger.Warning(fmt.Sprintf("Failed to tag %s with git tags, file must be unstaged", block.GetFilePath()))
	}
	var newTags []tags.ITag
	for _, tag := range t.Tags {
		tag, err := tag.CalculateValue(blame)
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to calculate tag value of tag %v, err: %s", tag, err))
			continue
		}
		newTags = append(newTags, tag)
	}
	block.AddNewTags(newTags)
}
