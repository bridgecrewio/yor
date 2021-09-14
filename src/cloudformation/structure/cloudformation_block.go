package structure

import (
	"reflect"

	"github.com/bridgecrewio/yor/src/common/structure"
)

type CloudformationBlock struct {
	structure.Block
}

func (b *CloudformationBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}

	mergedTags := b.MergeTags()
	cfnTags := make(map[string]string, len(mergedTags))
	for _, t := range mergedTags {
		cfnTags[t.GetKey()] = t.GetValue()
	}

	// set the tags attribute with the new tags
	reflect.ValueOf(b.RawBlock).Elem().FieldByName(b.TagsAttributeName).Set(reflect.ValueOf(cfnTags))
}

func (b *CloudformationBlock) GetTagsLines() structure.Lines {
	return b.TagLines
}

func (b *CloudformationBlock) GetSeparator() string {
	return "/n"
}
