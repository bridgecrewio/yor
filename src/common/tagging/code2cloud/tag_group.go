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

func (t *TagGroup) InitTagGroup(_ string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
	for _, fn := range options {
		fn(&t.Options)
	}
	t.SkippedTags = skippedTags
	t.SpecifiedTags = explicitlySpecifiedTags
	t.SetTags([]tags.ITag{&YorTraceTag{}, &YorNameTag{}})
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{
		&YorTraceTag{},
		&YorNameTag{},
	}
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return t.UpdateBlockTags(block)
}

func (t *TagGroup) UpdateBlockTags(block structure.IBlock) error {
	var newTags []tags.ITag
	var err error
	var tagVal tags.ITag
	for _, tag := range t.GetTags() {
		tagVal, err = tag.CalculateValue(block)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create %v tag for block %v", tag.GetKey(), block.GetResourceID()))
		}
		if tagVal != nil && tagVal.GetValue() != "" {
			newTags = append(newTags, tagVal)
		}
	}
	block.AddNewTags(newTags)
	return err
}
