package structure

import (
	"bridgecrewio/yor/common/tagging/tags"
)

type Block struct {
	FilePath          string
	ExitingTags       []tags.ITag
	NewTags           []tags.ITag
	RawBlock          interface{}
	IsTaggable        bool
	TagsAttributeName string
}

type IBlock interface {
	Init(filePath string, rawBlock interface{})
	String() string
	GetLines() []int
	GetExistingTags() []tags.ITag
	GetNewTags() []tags.ITag
	GetRawBlock() interface{}
	GetNewOwner() string
	GetPreviousOwner() string
	GetTraceId() string
	AddNewTags(newTags []tags.ITag)
	MergeTags() []tags.ITag
	CalculateTagsDiff() map[string][]tags.ITag
	IsBlockTaggable() bool
}

func (b *Block) AddNewTags(newTags []tags.ITag) {
	isTraced := false
	for _, tag := range b.ExitingTags {
		if _, ok := tag.(*tags.YorTraceTag); ok {
			isTraced = true
			break
		}
	}
	if isTraced {
		var yorTraceIndex int
		for index, tag := range newTags {
			if _, ok := tag.(*tags.YorTraceTag); ok {
				yorTraceIndex = index
			}
		}

		b.NewTags = append(b.NewTags[:yorTraceIndex], b.NewTags[yorTraceIndex+1:]...)
	}
	b.NewTags = append(b.NewTags, newTags...)
}

// MergeTags merges the tags and returns only the relevant Yor tags.
func (b *Block) MergeTags() []tags.ITag {
	var mergedTags []tags.ITag

	for _, tag := range b.ExitingTags {
		if val, ok := tag.(*tags.YorTraceTag); ok {
			mergedTags = append(mergedTags, val)
		}
	}

	mergedTags = append(mergedTags, b.NewTags...)

	return mergedTags
}

// CalculateTagsDiff returns a map which explains the changes in tags for this block
// added is the new tags, updated is the tags which were modified
func (b *Block) CalculateTagsDiff() map[string][]tags.ITag {
	var diff = make(map[string][]tags.ITag)
	for _, newTag := range b.GetNewTags() {
		found := false
		for _, existingTag := range b.GetExistingTags() {
			if newTag.GetKey() == existingTag.GetKey() {
				found = true
				if newTag.GetValue() != existingTag.GetValue() {
					diff["updated"] = append(diff["updated"], newTag)
					break
				}
			}
		}
		if !found {
			diff["added"] = append(diff["added"], newTag)
		}
	}
	return diff
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
