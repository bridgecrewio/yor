package structure

import (
	"reflect"

	"github.com/bridgecrewio/yor/src/common/structure"

	goformationTags "github.com/awslabs/goformation/v4/cloudformation/tags"
)

type CloudformationBlock struct {
	structure.Block
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

func (b *CloudformationBlock) GetTagsLines() structure.Lines {
	return b.TagLines
}

func (b *CloudformationBlock) GetSeparator() string {
	return "/n"
}
