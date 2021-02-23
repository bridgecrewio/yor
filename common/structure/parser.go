package structure

import "strings"

type IParser interface {
	Init(rootDir string, args map[string]string)
	ParseFile(filePath string) ([]IBlock, int, error)
	WriteFile(readFilePath string, blocks []IBlock, writeFilePath string) error
	GetSkippedDirs() []string
	GetAllowedFileTypes() []string
}

func IsFileSkipped(p IParser, file string) bool {
	matchingSuffix := false
	for _, suffix := range p.GetAllowedFileTypes() {
		if strings.HasSuffix(file, suffix) {
			matchingSuffix = true
		}
	}
	if !matchingSuffix {
		return true
	}
	for _, pattern := range p.GetSkippedDirs() {
		if strings.Contains(file, pattern) {
			return true
		}
	}
	return false
}
