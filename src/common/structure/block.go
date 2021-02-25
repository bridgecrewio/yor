package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/tagging/tags"
)

type Block struct {
	FilePath          string
	ExitingTags       []tags.ITag
	NewTags           []tags.ITag
	RawBlock          interface{}
	IsTaggable        bool
	TagsAttributeName string
}

type TagDiff struct {
	Added   []tags.ITag
	Updated []*tags.TagDiff
}

type IBlock interface {
	Init(filePath string, rawBlock interface{})
	String() string
	GetFilePath() string
	GetLines() []int
	GetExistingTags() []tags.ITag
	GetNewTags() []tags.ITag
	GetRawBlock() interface{}
	GetTraceID() string
	AddNewTags(newTags []tags.ITag)
	MergeTags() []tags.ITag
	CalculateTagsDiff() *TagDiff
	IsBlockTaggable() bool
	GetResourceID() string
}

func (b *Block) AddNewTags(newTags []tags.ITag) {
	isTraced := false
	var yorTag tags.YorTraceTag
	yorTag.Init()
	yorTagKeyName := yorTag.GetKey()
	for _, tag := range b.ExitingTags {
		match := common.IsTagKeyMatch(tag, yorTagKeyName)
		if _, ok := tag.(*tags.YorTraceTag); ok || match {
			isTraced = true
			break
		}
	}
	if isTraced {
		var yorTraceIndex int
		for index, tag := range newTags {
			match := common.IsTagKeyMatch(tag, yorTagKeyName)
			if _, ok := tag.(*tags.YorTraceTag); ok || match {
				yorTraceIndex = index
			}
		}

		newTags = append(newTags[:yorTraceIndex], newTags[yorTraceIndex+1:]...)
	}
	b.NewTags = append(b.NewTags, newTags...)
}

// MergeTags merges the tags and returns only the relevant Yor tags.
func (b *Block) MergeTags() []tags.ITag {
	var mergedTags []tags.ITag
	var yorTag tags.YorTraceTag
	yorTag.Init()
	yorTagKeyName := yorTag.GetKey()
	for _, tag := range b.ExitingTags {
		match := common.IsTagKeyMatch(tag, yorTagKeyName)
		if val, ok := tag.(*tags.YorTraceTag); ok || match {
			if val != nil {
				mergedTags = append(mergedTags, val)
			} else {
				mergedTags = append(mergedTags, tag)
			}
		}
	}

	mergedTags = append(mergedTags, b.NewTags...)

	return mergedTags
}

// CalculateTagsDiff returns a map which explains the changes in tags for this block
// Added is the new tags, Updated is the tags which were modified
func (b *Block) CalculateTagsDiff() *TagDiff {
	var diff = TagDiff{}
	for _, newTag := range b.GetNewTags() {
		found := false
		for _, existingTag := range b.GetExistingTags() {
			if newTag.GetKey() == existingTag.GetKey() {
				found = true
				if newTag.GetValue() != existingTag.GetValue() {
					diff.Updated = append(diff.Updated, &tags.TagDiff{
						Key:       newTag.GetKey(),
						PrevValue: existingTag.GetValue(),
						NewValue:  newTag.GetValue(),
					})
					break
				}
			}
		}
		if !found {
			diff.Added = append(diff.Added, newTag)
		}
	}
	return &diff
}

func (b *Block) GetRawBlock() interface{} {
	return b.RawBlock
}

func (b *Block) GetExistingTags() []tags.ITag {
	return b.ExitingTags
}

func (b *Block) GetNewTags() []tags.ITag {
	return b.NewTags
}

func (b *Block) IsBlockTaggable() bool {
	return b.IsTaggable
}

func (b *Block) GetFilePath() string {
	return b.FilePath
}

func (b *Block) GetTraceID() string {
	for _, tag := range b.GetExistingTags() {
		if val, ok := tag.(*tags.YorTraceTag); ok {
			return val.GetValue()
		}
	}
	for _, tag := range b.GetNewTags() {
		if val, ok := tag.(*tags.YorTraceTag); ok {
			return val.GetValue()
		}
	}
	return ""
}
