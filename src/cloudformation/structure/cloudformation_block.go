package structure

import (
	"bridgecrewio/yor/src/common"
	"reflect"

	goformationTags "github.com/awslabs/goformation/v4/cloudformation/tags"
)

type CloudformationBlock struct {
	common.Block
	Name string
}

func (b *CloudformationBlock) GetResourceID() string {
	return b.Name
}

func (b *CloudformationBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *CloudformationBlock) GetLines(_ ...bool) common.Lines {
	return b.Lines
}

func (b *CloudformationBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}

	mergedTags := b.MergeTags()
	cfnMergedTags := make([]goformationTags.Tag, 0)
	for _, t := range mergedTags {
		cfnMergedTags = append(cfnMergedTags, goformationTags.Tag{
			Key:   t.GetKey(),
			Value: t.GetValue(),
		})
	}

	// set the tags attribute with the new tags
	reflect.ValueOf(b.RawBlock).Elem().FieldByName(b.TagsAttributeName).Set(reflect.ValueOf(cfnMergedTags))
}

func (b *CloudformationBlock) GetTagsLines() common.Lines {
	return b.TagLines
}

func (b *CloudformationBlock) GetSeparator() string {
	return "/n"
}
