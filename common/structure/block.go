package structure

import (
	"bridgecrewio/yor/common/tagging/tags"
)

type Block struct {
	FilePath    string
	ExitingTags []tags.ITag
	NewTags     []tags.ITag
	RawBlock    interface{}
}

type IBlock interface {
	Init(filePath string, rawBlock interface{})
	String() string
	GetLines() []int
	GetRawBlock() interface{}
	GetNewOwner() string
	GetPreviousOwner() string
	GetTraceId() string
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
	// TODO
	return nil
}
