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
	b.NewTags = append(b.NewTags, newTags...)
}

// MergeTags merges the tags while trying to preserve the existing order of the tags.
// It does this by first iterating through the existing tags and updating the values if there's anything to update.
// Then it adds all the new tags that didn't exist beforehand.
func (b *Block) MergeTags() []tags.ITag {
	var mergedTags []tags.ITag

	for _, existingTag := range b.GetExistingTags() {
		found := false
		for _, newTag := range b.GetNewTags() {
			if newTag.GetKey() == existingTag.GetKey() {
				mergedTags = append(mergedTags, newTag)
				found = true
				break
			}
		}
		if !found {
			mergedTags = append(mergedTags, existingTag)
		}
	}
	for _, newTag := range b.GetNewTags() {
		var found bool
		for _, updateTag := range mergedTags {
			if updateTag.GetKey() == newTag.GetKey() {
				found = true
				break
			}
		}
		if !found {
			mergedTags = append(mergedTags, newTag)
		}
	}
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
				if newTag.GetValue() != existingTag.GetValue() {
					diff["updated"] = append(diff["updated"], newTag)
					found = true
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
