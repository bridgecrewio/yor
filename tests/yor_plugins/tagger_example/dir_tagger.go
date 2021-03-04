package main

import (
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
)

type DirTagger struct {
	tagging.Tagger
}

func (d *DirTagger) InitTagger(_ string, skippedTags []string) {
	// If skipped tags isn't passed in, the skip mechanism will not work
	d.SkippedTags = skippedTags
	d.SetTags([]tags.ITag{&DirTag{}})
}

func (d *DirTagger) CreateTagsForBlock(block structure.IBlock) {
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
