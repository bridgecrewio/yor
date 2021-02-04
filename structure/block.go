package structure

type Block struct {
	FilePath    string
	ExitingTags map[string]interface{}
	NewTags     map[string]interface{}
	RawBlock    interface{}
}

type iBlock interface {
	Init(filePath string, rawBlock interface{})
	String() string
	GetLines() []int
	GetRawBlock() interface{}
}

func (b *Block) AddNewTags(newTags map[string]interface{}) {
	// TODO
}

func (b *Block) MergeTags() map[string]interface{} {
	// TODO - return a map of the old and new tags
	return nil
}

func (b *Block) CalculateTagsDiff() map[string][]map[string]interface{} {
	// TODO - return a map with keys such as "added", "deleted", modified" and the matching tags
	return nil
}

func (b *Block) GetRawBlock() interface{} {
	// TODO
	return nil
}
