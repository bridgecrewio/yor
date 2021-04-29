package common

import "bridgecrewio/yor/src/common/structure"

type IParser interface {
	Init(rootDir string, args map[string]string)
	ParseFile(filePath string) ([]structure.IBlock, error)
	WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error
	GetSkippedDirs() []string
	GetSupportedFileExtensions() []string
}
