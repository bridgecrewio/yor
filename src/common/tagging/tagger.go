package tagging

import (
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
)

type Tagger struct {
	Tags []tags.ITag
}

var IgnoredDirs = []string{".git", ".DS_Store", ".idea"}

type ITagger interface {
	InitTagger(path string)
	InitTags(extraTags []tags.ITag)
	CreateTagsForBlock(block structure.IBlock)
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
