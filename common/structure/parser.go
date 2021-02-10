package structure

type IParser interface {
	Init()
	ParseFile(filePath string, rootDir string) ([]IBlock, error)
	WriteFile(filePath string, blocks []IBlock) error
}
