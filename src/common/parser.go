package common

import "github.com/bridgecrewio/yor/src/common/structure"

/*
	This interface is for all parsers.  All these signatures need to be implemented
*/
type IParser interface {
	Init(rootDir string, args map[string]string)
	Name() string
	ValidFile(filePath string) bool
	ParseFile(filePath string) ([]structure.IBlock, error)
	WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error
	GetSkippedDirs() []string
	GetSupportedFileExtensions() []string
	Close()
}
