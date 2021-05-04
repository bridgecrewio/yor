package structure

import (
	"github.com/bridgecrewio/yor/src/common/structure"

	"go.opencensus.io/tag"
)

type ServerlessBlock struct {
	structure.Block
	Name string
}

func (b *ServerlessBlock) GetResourceID() string {
	return b.Name
}

func (b *ServerlessBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *ServerlessBlock) GetLines(_ ...bool) structure.Lines {
	return b.Block.Lines
}

func (b *ServerlessBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}

	mergedTags := b.MergeTags()
	slsMergedTags := make([]tag.Tag, 0)
	for _, t := range mergedTags {
		key, _ := tag.NewKey(t.GetKey())
		slsTag := tag.Tag{
			Key:   key,
			Value: t.GetValue(),
		}

		slsMergedTags = append(slsMergedTags, slsTag)
	}
	slsMergedTagsValue := make(map[string]string)
	for _, mergedTag := range slsMergedTags {
		slsMergedTagsValue[mergedTag.Key.Name()] = mergedTag.Value
	}
	b.RawBlock.(map[interface{}]interface{})[b.TagsAttributeName] = slsMergedTagsValue
}

func (b *ServerlessBlock) GetTagsLines() structure.Lines {
	return b.Block.TagLines
}

func (b *ServerlessBlock) GetSeparator() string {
	return "/n"
}
