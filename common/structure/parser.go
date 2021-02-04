package structure

type IParser interface {
	ParseFile(filePath string) ([]*Block, error)
	WriteFile(filePath string, blocks []*Block) error
}
