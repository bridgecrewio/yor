package main

import (
	"fmt"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type OrgTagGroup struct {
	tagging.TagGroup
}

func (d *OrgTagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{
		&DirTag{},
	}
}

func (d *OrgTagGroup) InitTagGroup(_ string, skippedTags []string) {
	// If skipped tags isn't passed in, the skip mechanism will not work
	d.SkippedTags = skippedTags
	d.SetTags(d.GetDefaultTags())
}

func (d *OrgTagGroup) CreateTagsForBlock(block structure.IBlock) {
	var newTags []tags.ITag
	for _, tag := range d.GetTags() {
		tagVal, err := tag.CalculateValue(block)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create %v tag for block %v", tag.GetKey(), block.GetResourceID()))
		}
		newTags = append(newTags, tagVal)
	}
	block.AddNewTags(newTags)
}
