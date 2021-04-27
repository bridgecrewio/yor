package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/structure"
	"go.opencensus.io/tag"
	"reflect"
)

type ServerlessBlock struct {
	structure.Block
	name string
}

func (b *ServerlessBlock) GetResourceID() string {
	return b.name
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
	slsMergedTagsValue := make(map[reflect.Value]reflect.Value, 0)
	for _, mergedTag := range slsMergedTags {
		slsMergedTagsValue[reflect.ValueOf(mergedTag.Key.Name())] = reflect.ValueOf(mergedTag.Value)
	}
	var someMap map[interface{}]interface{}
	tagsValueRef := reflect.MakeMap(reflect.TypeOf(someMap))
	for i, mapKey := range slsMergedTags {
		tagsValueRef.SetMapIndex(reflect.ValueOf(mapKey.Key.Name()), reflect.ValueOf(slsMergedTags[i].Value))
	}
	b.RawBlock.(reflect.Value).SetMapIndex(reflect.ValueOf(b.TagsAttributeName), reflect.ValueOf(tagsValueRef))
}

func (b *ServerlessBlock) GetTagsLines() common.Lines {
	return b.Block.TagLines
}

func (b *ServerlessBlock) GetSeparator() string {
	return "/n"
}
