package structure

import (
	"github.com/bridgecrewio/yor/src/common/structure"
	"go.opencensus.io/tag"
)

type ServerlessBlock struct {
	structure.Block
}

func (b *ServerlessBlock) GetFramework() string {
	return "Serverless"
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
	slsMergedTagsValue := make(map[string]interface{})
	for _, mergedTag := range slsMergedTags {
		slsMergedTagsValue[mergedTag.Key.Name()] = mergedTag.Value
	}
	rawFunction := b.RawBlock.(structure.Function)
	if rawFunction.Tags == nil {
		rawFunction.Tags = make(map[string]interface{}, len(slsMergedTags))
	}
	for _, slsTag := range slsMergedTags {
		rawFunction.Tags[slsTag.Key.Name()] = slsTag.Value
	}
	b.RawBlock = rawFunction
}

func (b *ServerlessBlock) GetTagsLines() structure.Lines {
	return b.Block.TagLines
}

func (b *ServerlessBlock) GetSeparator() string {
	return "/n"
}
