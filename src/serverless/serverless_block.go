package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/structure"
	"fmt"
	"go.opencensus.io/tag"
	"reflect"
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
	slsMergedTags := make([]tag.Tag, 0)
	for _, t := range mergedTags {
		key, _ := tag.NewKey(t.GetKey())
		slsTag := tag.Tag{
			Key:   key,
			Value: t.GetValue(),
		}

		slsMergedTags = append(slsMergedTags, slsTag)
	}
	slsMergedTagsValue := make([]reflect.Value, 0)
	for _, mergedTag := range slsMergedTags {
		slsMergedTagsValue = append(slsMergedTagsValue, reflect.ValueOf(map[string]string{mergedTag.Key.Name(): mergedTag.Value}))
	}
	for numField := 0; numField < reflect.ValueOf(b.RawBlock).NumField(); numField++ {
		field := reflect.ValueOf(b.RawBlock).Field(numField).Elem()
		if field.Kind() == reflect.Struct {
			for numSubField := 0; numSubField < reflect.ValueOf(field).NumField(); numSubField++ {
				subField := field.Field(numSubField)
				fmt.Println(subField.Elem().Kind())
			}
		}
		fmt.Println(field)
		//if field == b.TagsAttributeName {
		//	fmt.Println(1)
		//	//blockValue.Set(reflect.ValueOf(slsMergedTagsValue))
		//	break
		//}
	}
	for _, blockValue := range reflect.ValueOf(b.RawBlock).MapKeys() {
		if blockValue.Elem().String() == b.TagsAttributeName {
			for _, a := range blockValue.Elem().MapKeys() {
				fmt.Println(a)
			}
		}
	}
	//set the tags attribute with the new tags
	reflect.ValueOf(b.RawBlock).SetMapIndex(reflect.ValueOf(b.TagsAttributeName), reflect.ValueOf(slsMergedTagsValue))
	fmt.Println(b)
}

func (b *ServerlessBlock) GetTagsLines() common.Lines {
	return b.tagLines
}

func (b *ServerlessBlock) GetSeparator() string {
	return "/n"
}
