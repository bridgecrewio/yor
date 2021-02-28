package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/code2cloud"
	"bridgecrewio/yor/src/common/tagging/tags"
	"reflect"

	goformation_tags "github.com/awslabs/goformation/v4/cloudformation/tags"
)

type CloudformationBlock struct {
	structure.Block
	lines common.Lines
	name  string
}

func (b *CloudformationBlock) GetResourceID() string {
	return b.name
}

func (b *CloudformationBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *CloudformationBlock) GetLines() common.Lines {
	return b.lines
}

func (b *CloudformationBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}

	mergedTags := b.MergeCFNTags()
	cfnMergedTags := make([]goformation_tags.Tag, 0)
	for _, t := range mergedTags {
		cfnMergedTags = append(cfnMergedTags, goformation_tags.Tag{
			Key:   t.GetKey(),
			Value: t.GetValue(),
		})
	}

	// set the tags attribute with the new tags
	reflect.ValueOf(b.RawBlock).Elem().FieldByName(b.TagsAttributeName).Set(reflect.ValueOf(cfnMergedTags))
}

func (b *CloudformationBlock) MergeCFNTags() []tags.ITag {
	existingTagsByKey := map[string]tags.ITag{}
	newTagsByKey := map[string]tags.ITag{}

	for _, tag := range b.ExitingTags {
		existingTagsByKey[tag.GetKey()] = tag
	}
	for _, tag := range b.NewTags {
		newTagsByKey[tag.GetKey()] = tag
	}

	var mergedTags []tags.ITag
	var yorTag code2cloud.YorTraceTag
	yorTag.Init()
	yorTagKeyName := yorTag.GetKey()
	for _, existingTag := range b.ExitingTags {
		if newTag, ok := newTagsByKey[existingTag.GetKey()]; ok {
			match := tags.IsTagKeyMatch(existingTag, yorTagKeyName)
			if val, ok := existingTag.(*code2cloud.YorTraceTag); ok || match {
				if val != nil {
					mergedTags = append(mergedTags, val)
				} else {
					mergedTags = append(mergedTags, existingTag)
				}
			} else {
				mergedTags = append(mergedTags, newTag)
			}
			delete(newTagsByKey, existingTag.GetKey())
		} else {
			mergedTags = append(mergedTags, existingTag)
		}
	}

	for newTagKey := range newTagsByKey {
		mergedTags = append(mergedTags, newTagsByKey[newTagKey])
	}

	return mergedTags
}
