package code2cloud

import (
	"fmt"

	"github.com/bridgecrewio/yor/src/common/logger"
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

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) {
	var newTags []tags.ITag
	for _, tag := range t.GetTags() {
		tagVal, err := tag.CalculateValue(struct{}{})
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create %v tag for block %v", tag.GetKey(), block.GetResourceID()))
		}
		newTags = append(newTags, tagVal)
	}
	block.AddNewTags(newTags)
}
