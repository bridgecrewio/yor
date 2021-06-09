package common

import "github.com/bridgecrewio/yor/src/common/structure"

type IParser interface {
	Init(rootDir string, args map[string]string)
	Name() string
	ParseFile(filePath string) ([]structure.IBlock, error)
	WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error
	GetSkippedDirs() []string
	GetSupportedFileExtensions() []string
}
