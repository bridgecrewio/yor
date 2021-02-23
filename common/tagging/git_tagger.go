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
	GitService *gitservice.GitService
}

func (t *GitTagger) InitTagger(path string) {
	gitService, err := gitservice.NewGitService(path)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize git service for path %s", path))
	}
	t.GitService = gitService
}

func (t *GitTagger) TagFile(path string, fileLength int) bool {
	fileBlame, err := t.GitService.GetFileBlame(path)
	if err != nil {
		logger.Warning(fmt.Sprintf("Unable to get git blame for file %s: %s", path, err))
		return false
	}
	if len(fileBlame.Lines) != fileLength {
		logger.Warning(fmt.Sprintf("Unable to tag file %s because the file contains uncommitted changes", path))
		return false
	}

	return true
}

func (t *GitTagger) CreateTagsForBlock(block structure.IBlock) {
	blame, err := t.GitService.GetBlameForFileLines(block.GetFilePath(), block.GetLines())
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to tag %v with git tags, err: %v", block.GetResourceID(), err.Error()))
		return
	}
	if blame == nil {
		logger.Warning(fmt.Sprintf("Failed to tag %s with git tags, file must be unstaged", block.GetFilePath()))
	}
	var newTags []tags.ITag
	for _, tag := range t.Tags {
		newTag, err := tag.CalculateValue(blame)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to calculate tag value of tag %v, err: %s", tag.GetKey(), err))
			continue
		}
		newTags = append(newTags, newTag)
	}
	block.AddNewTags(newTags)
}
