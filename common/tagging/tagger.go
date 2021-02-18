package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
)

type Tagger struct {
	Tags []tags.ITag
}

var IgnoredDirs = []string{".git", ".DS_Store", ".idea"}

type ITagger interface {
	InitTags(extraTags []tags.ITag)
	CreateTagsForBlock(block structure.IBlock, gitBlame *gitservice.GitBlame)
	IsFileSkipped(file string) bool
}

func (t *Tagger) InitTags(extraTags []tags.ITag) {
	for _, tagType := range tags.TagTypes {
		tagType.Init()
	}
	t.Tags = append(t.Tags, tags.TagTypes...)
	t.Tags = append(t.Tags, extraTags...)
}

func (t *Tagger) GetSkippedDirs() []string {
	return IgnoredDirs
}
