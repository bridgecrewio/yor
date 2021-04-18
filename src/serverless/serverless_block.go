package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/structure"
	"reflect"

	goformation_tags "github.com/awslabs/goformation/v4/cloudformation/tags"
)

type ServerlessBlock struct {
	structure.Block
	lines    common.Lines
	name     string
	tagLines common.Lines
}

func (b *ServerlessBlock) GetResourceID() string {
	return b.name
}

func (b *ServerlessBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *ServerlessBlock) GetLines(_ ...bool) common.Lines {
	return b.lines
}

func (b *ServerlessBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}

	mergedTags := b.MergeTags()
	slsMergedTags := make([]goformation_tags.Tag, 0)
	for _, t := range mergedTags {
		slsMergedTags = append(slsMergedTags, goformation_tags.Tag{
			Key:   t.GetKey(),
			Value: t.GetValue(),
		})
	}

	// set the tags attribute with the new tags
	reflect.ValueOf(b.RawBlock).Elem().FieldByName(b.TagsAttributeName).Set(reflect.ValueOf(slsMergedTags))
}

func (b *ServerlessBlock) GetTagsLines() common.Lines {
	return b.tagLines
}

func (b *ServerlessBlock) GetSeparator() string {
	return "/n"
}
