package structure

type IParser interface {
	Init(rootDir string, args map[string]string)
	ParseFile(filePath string) ([]IBlock, error)
	WriteFile(filePath string, blocks []IBlock) error
}
