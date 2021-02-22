package structure

import "strings"

type IParser interface {
	Init(rootDir string, args map[string]string)
	ParseFile(filePath string) ([]IBlock, error)
	WriteFile(readFilePath string, blocks []IBlock, writeFilePath string) error
	IsFileSkipped(file string) bool
	GetSkippedDirs() []string
}

type Parser struct{}

func (p *Parser) IsFileSkipped(file string) bool {
	for _, pattern := range p.GetSkippedDirs() {
		if strings.Contains(file, pattern) {
			return true
		}
	}
	return false
}

func (p *Parser) Init(rootDir string, args map[string]string) {
	panic("implement me")
}

func (p *Parser) ParseFile(filePath string) ([]IBlock, error) {
	panic("implement me")
}

func (p *Parser) WriteFile(readFilePath string, blocks []IBlock, writeFilePath string) error {
	panic("implement me")
}

func (p *Parser) GetSkippedDirs() []string {
	return []string{}
}
