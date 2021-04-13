package structure

type IParser interface {
	Init(rootDir string, args map[string]string)
	ParseFile(filePath string) ([]IBlock, error)
	WriteFile(readFilePath string, blocks []IBlock, writeFilePath string) error
	GetSkippedDirs() []string
	GetSupportedFileExtensions() []string
}
