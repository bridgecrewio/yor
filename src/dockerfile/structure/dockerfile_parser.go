package structure

import (
	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/minamijoyo/tfschema/tfschema"
	"sync"
)

type DockerfileParser struct {
	Labels                 map[string]string
	rootDir                string
	taggableResourcesCache map[string]bool
	providerToClientMap    sync.Map
	dockerFileExists       bool
	NewLabels              []TagItem
}

func (p *DockerfileParser) Name() string {
	return "Dockerfile"
}

func (p *DockerfileParser) ValidFile(_ string) bool {
	return true
}

func (p *DockerfileParser) ParseFile(filePath string) ([]structure.IBlock, error) {

	parsedBlocks := make([]structure.IBlock, 0)
	return parsedBlocks, nil

}

func (p *DockerfileParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	return nil
}

func (p *DockerfileParser) GetSkippedDirs() []string {
	var ignoredDirs []string
	return ignoredDirs
}

func (p *DockerfileParser) Close() {
	logger.MuteOutputBlock(func() {
		p.providerToClientMap.Range(func(provider, iClient interface{}) bool {
			client := iClient.(tfschema.Client)
			client.Close()
			return true
		})
	})
}

func (p *DockerfileParser) GetSupportedFileExtensions() []string {
	return []string{common.TfFileType.Extension}
}

func (p *DockerfileParser) Init(rootDir string, args map[string]string) {
	p.rootDir = rootDir
	p.dockerFileExists = false
	p.IngestFile()
	p.ReadYamlConfiguration()
	p.PatchDockerFile()

}

type TagItem struct {
	Key           string
	Value         string
	NeedsUpdating bool
	Exists        bool
}
