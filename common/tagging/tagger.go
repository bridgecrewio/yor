package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
)

type Tagger struct {
	Tags []tags.ITag
}

type ITagger interface {
	CreateTagsForBlock(block structure.IBlock, gitBlame *git_service.GitBlame) error // this method should call IBlock.AddNewTags
	InitTags(extraTags []tags.ITag)
}

func (t *Tagger) InitTags(extraTags []tags.ITag) {
	for _, tagType := range tags.TagTypes {
		t.Tags = append(t.Tags, tagType.Init())
	}

	for _, extraTag := range extraTags {
		t.Tags = append(t.Tags, extraTag)
	}
}
