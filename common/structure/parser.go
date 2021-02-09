package structure

type IParser interface {
	ParseFile(filePath string) ([]IBlock, error)
	WriteFile(filePath string, blocks []IBlock) error
}
