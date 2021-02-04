package structure

import "bridgecrewio/yor/common/tagging"

type Block struct {
	FilePath    string
	ExitingTags []tagging.ITag
	NewTags     []tagging.ITag
	RawBlock    interface{}
}

type IBlock interface {
	Init(filePath string, rawBlock interface{})
	String() string
	GetLines() []int
	GetRawBlock() interface{}
}

func (b *Block) AddNewTags(newTags []tagging.ITag) {
	// TODO
}

func (b *Block) MergeTags() []tagging.ITag {
	// TODO - return a map of the old and new tags
	return nil
}

func (b *Block) CalculateTagsDiff() map[string][]tagging.ITag {
	// TODO - return a map with keys such as "added", "deleted", modified" and the matching tags
	return nil
}

func (b *Block) GetRawBlock() interface{} {
	// TODO
	return nil
}
