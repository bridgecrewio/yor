package structure

type IParser interface {
	Init(rootDir string)
	ParseFile(filePath string) ([]IBlock, error)
	WriteFile(filePath string, blocks []IBlock) error
}
