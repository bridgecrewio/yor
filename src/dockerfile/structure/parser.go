package structure

import (
	"fmt"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"os"
	"strings"
)

type DockerfileParser struct {
}

func (p *DockerfileParser) ValidFile(filePath string) bool {
	fileNameParts := strings.Split(filePath, "/")
	fileName := fileNameParts[len(fileNameParts)-1]
	return strings.HasPrefix(fileName, "Dockerfile") || strings.HasPrefix(fileName, "dockerfile") ||
		strings.HasSuffix(fileName, "Dockerfile") || strings.HasSuffix(fileName, "dockerfile")
}

func (p *DockerfileParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	data, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("readfile error: %w", err)
	}

	result, err := parser.Parse(data)

	_ = result

	defer func(data *os.File) {
		err := data.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("close error:%s", err))
		}
	}(data)

	return nil, err
}

func (p *DockerfileParser) Init(rootDir string, _ map[string]string) {
	//TODO implement me
}

func (p *DockerfileParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	//TODO implement me
	panic("implement me")
}

func (p *DockerfileParser) GetSkippedDirs() []string {
	//TODO implement me
	panic("implement me")
}

func (p *DockerfileParser) GetSupportedFileExtensions() []string {
	//TODO implement me
	panic("implement me")
}

func (p *DockerfileParser) Close() {
	//TODO implement me
	panic("implement me")
}

func (p *DockerfileParser) Name() string {
	return "Dockerfile"
}
