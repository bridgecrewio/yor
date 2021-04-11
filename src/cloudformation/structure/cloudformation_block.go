package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/structure"
	"reflect"

	goformation_tags "github.com/awslabs/goformation/v4/cloudformation/tags"
)

type CloudformationBlock struct {
	structure.Block
	lines    common.Lines
	name     string
	tagLines common.Lines
}

func (b *CloudformationBlock) GetResourceID() string {
	return b.name
}

func (b *CloudformationBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *CloudformationBlock) GetLines(_ ...bool) common.Lines {
	return b.lines
}

func (b *CloudformationBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}

	mergedTags := b.MergeTags()
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

func (b *CloudformationBlock) GetTagsLines() common.Lines {
	return b.lines
}

func (b *CloudformationBlock) GetSeparator() string {
	return "/n"
}
