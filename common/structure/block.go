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
	// TODO
}

func (b *Block) MergeTags() []tags.ITag {
	// TODO - return a map of the old and new tags
	return nil
}

func (b *Block) CalculateTagsDiff() map[string][]tags.ITag {
	// TODO - return a map with keys such as "added", "deleted", modified" and the matching tags
	return nil
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
