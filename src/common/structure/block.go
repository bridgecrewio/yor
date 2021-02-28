package structure

import (
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
	if newTags == nil {
		return
	}
	var tagsToAdd []tags.ITag
	traceTagKey := tags.YorTraceTagKey
	for _, newTag := range newTags {
		found := false
		for _, existingTag := range b.ExitingTags {
			if existingTag.GetKey() == newTag.GetKey() {
				found = true
				if existingTag.GetKey() == traceTagKey || existingTag.GetValue() == newTag.GetValue() {
					continue
				}
				tagsToAdd = append(tagsToAdd, newTag)
			}
		}
		if !found {
			tagsToAdd = append(tagsToAdd, newTag)
		}
	}

	b.NewTags = append(b.NewTags, tagsToAdd...)
}

// MergeTags merges the tags and returns only the relevant Yor tags.
func (b *Block) MergeTags() []tags.ITag {
	var mergedTags []tags.ITag
	yorTagKeyName := tags.YorTraceTagKey
	for _, tag := range b.ExitingTags {
		match := tags.IsTagKeyMatch(tag, yorTagKeyName)
		if match {
			mergedTags = append(mergedTags, tag)
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
	for _, tag := range b.MergeTags() {
		if tag.GetKey() == tags.YorTraceTagKey {
			return tag.GetValue()
		}
	}
	return ""
}
