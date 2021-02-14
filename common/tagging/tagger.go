package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
)

type Tagger struct {
	Tags []tags.ITag
}

type ITagger interface {
	InitTags(extraTags []tags.ITag)
	CreateTagsForBlock(block structure.IBlock, gitBlame *gitservice.GitBlame)
}

func (t *Tagger) InitTags(extraTags []tags.ITag) {
	for _, tagType := range tags.TagTypes {
		tagType.Init()
		t.Tags = append(t.Tags, tagType)
	}
	t.Tags = append(t.Tags, extraTags...)
}
