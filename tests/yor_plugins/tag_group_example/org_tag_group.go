package main

import (
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type OrgTagGroup struct {
	tagging.TagGroup
}

func (d *OrgTagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return d.UpdateBlockTags(block, block)
}

func (d *OrgTagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{
		&DirTag{},
	}
}

func (d *OrgTagGroup) InitTagGroup(_ string, skippedTags []string, explicitlySpecifiedTags []string, tagPrefix string) {
	// If skipped tags isn't passed in, the skip mechanism will not work
	d.SkippedTags = skippedTags
	d.SpecifiedTags = explicitlySpecifiedTags
	d.SetTags(d.GetDefaultTags(), tagPrefix)
}
