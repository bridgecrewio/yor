package structure

import (
	"bridgecrewio/yor/src/common"
	"go.opencensus.io/tag"
)

type ServerlessBlock struct {
	common.Block
	Name string
}

func (b *ServerlessBlock) GetResourceID() string {
	return b.Name
}

func (b *ServerlessBlock) Init(filePath string, rawBlock interface{}) {
	b.RawBlock = rawBlock
	b.FilePath = filePath
}

func (b *ServerlessBlock) GetLines(_ ...bool) common.Lines {
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
	slsMergedTagsValue := make(map[string]string, 0)
	for _, mergedTag := range slsMergedTags {
		slsMergedTagsValue[mergedTag.Key.Name()] = mergedTag.Value
	}
	b.RawBlock.(map[interface{}]interface{})[b.TagsAttributeName] = slsMergedTagsValue
}

func (b *ServerlessBlock) GetTagsLines() common.Lines {
	return b.Block.TagLines
}

func (b *ServerlessBlock) GetSeparator() string {
	return "/n"
}
