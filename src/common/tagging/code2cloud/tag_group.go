package code2cloud

import (
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type TagGroup struct {
	tagging.TagGroup
}

func (t *TagGroup) InitTagGroup(_ string, skippedTags []string) {
	t.SkippedTags = skippedTags
	t.SetTags([]tags.ITag{&YorTraceTag{}})
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{
		&YorTraceTag{},
	}
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return t.UpdateBlockTags(block, struct{}{})
}
