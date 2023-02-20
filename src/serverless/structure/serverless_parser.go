package structure

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/types"
	"github.com/bridgecrewio/yor/src/common/utils"
	yamlUtils "github.com/bridgecrewio/yor/src/common/yaml"
	"github.com/thepauleh/goserverless"
	"github.com/thepauleh/goserverless/serverless"
)

const FunctionTagsAttributeName = "tags"
const FunctionsSectionName = "functions"

type ServerlessParser struct {
	YamlParser types.YamlParser
}

func (p *ServerlessParser) Name() string {
	return "Serverless"
}

func (p *ServerlessParser) Init(rootDir string, _ map[string]string) {
	p.YamlParser.RootDir = rootDir
}

func (p *ServerlessParser) Close() {
	return
}

func (p *ServerlessParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *ServerlessParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension}
}

func goserverlessParse(file string) (*serverless.Template, error) {
	var template *serverless.Template
	var err error
	defer func() {
		if e := recover(); e != nil {
			logger.Warning(fmt.Sprintf("Failed to parser serverless yaml at %v due to: %v", file, e))
			err = fmt.Errorf("failed to parse sls file %v: %v", file, e)
		}
	}()

	template, err = goserverless.Open(file)
	return template, err
}

func (p *ServerlessParser) ValidFile(_ string) bool {
	return true
}

func (p *ServerlessParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	parsedBlocks := make([]structure.IBlock, 0)
	fileFormat := utils.GetFileFormat(filePath)
	fileName := filepath.Base(filePath)
	if !(fileName == fmt.Sprintf("serverless.%s", fileFormat)) {
		return nil, nil
	}
	// #nosec G304 - file is from user
	template, err := goserverlessParse(filePath)
	if err != nil || template == nil || template.Functions == nil {
		if err != nil {
			logger.Warning(fmt.Sprintf("There was an error processing the serverless template: %s", err))
		}
		if err == nil {
			err = fmt.Errorf("failed to parse file %v", filePath)
		}
		return nil, err
	}
	if template.Functions == nil && template.Resources.Resources == nil {
		return parsedBlocks, nil
	}

	// cfnStackTagsResource := p.template.Provider.CFNTags
	resourceNames := make([]string, 0)
	var resourceNamesToLines map[string]*structure.Lines
	for funcName := range template.Functions {
		resourceNames = append(resourceNames, funcName)
	}
	switch utils.GetFileFormat(filePath) {
	case common.YmlFileType.FileFormat, common.YamlFileType.FileFormat:
		resourceNamesToLines = yamlUtils.MapResourcesLineYAML(filePath, resourceNames, FunctionsSectionName)
	default:
		return nil, fmt.Errorf("unsupported file type %s", utils.GetFileFormat(filePath))
	}
	minResourceLine := math.MaxInt8
	maxResourceLine := 0
	for _, funcName := range resourceNames {
		var existingTags []tags.ITag
		var slsBlock *ServerlessBlock
		tagsLines := structure.Lines{Start: -1, End: -1}
		var lines *structure.Lines
		slsFunction := template.Functions[funcName]
		lines = resourceNamesToLines[funcName]
		minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
		maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))
		if slsFunction.Tags != nil {
			tagsLines = p.getTagsLines(filePath, lines)
			for tagKey, tagValue := range slsFunction.Tags {
				existingTags = append(existingTags, &tags.Tag{Key: tagKey, Value: fmt.Sprintf("%v", tagValue)})
			}
		}

		slsBlock = &ServerlessBlock{
			Block: structure.Block{
				FilePath:          filePath,
				ExitingTags:       existingTags,
				RawBlock:          slsFunction,
				IsTaggable:        true,
				TagsAttributeName: FunctionTagsAttributeName,
				Lines:             *lines,
				TagLines:          tagsLines,
				Name:              funcName,
			},
		}

		parsedBlocks = append(parsedBlocks, slsBlock)
		p.YamlParser.FileToResourcesLines.Store(filePath, structure.Lines{Start: minResourceLine, End: maxResourceLine})

	}
	return parsedBlocks, nil
}

func (p *ServerlessParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, addToggle bool) error {
	for _, block := range blocks {
		block := block.(*ServerlessBlock)
		block.UpdateTags()
	}
	tempFile, err := os.CreateTemp(filepath.Dir(readFilePath), "temp.*.yaml")
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()
	if err != nil {
		return err
	}
	err = yamlUtils.WriteYAMLFile(readFilePath, blocks, tempFile.Name(), FunctionTagsAttributeName, FunctionsSectionName)
	if err != nil {
		return err
	}
	_, err = p.ParseFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("editing file %v resulted in a malformed template, please open a github issue with the relevant details", readFilePath)
	}
	return yamlUtils.WriteYAMLFile(readFilePath, blocks, writeFilePath, FunctionTagsAttributeName, FunctionsSectionName)
}

func (p *ServerlessParser) getTagsLines(filePath string, resourceLinesRange *structure.Lines) structure.Lines {
	nonFoundLines := structure.Lines{Start: -1, End: -1}
	fileFormat := utils.GetFileFormat(filePath)
	tagsLines := structure.Lines{Start: -1, End: -1}
	lineCounter := 0
	switch fileFormat {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		scanner, _ := utils.GetFileScanner(filePath, &nonFoundLines)
		// iterate file line by line
		tagsIndentSize := 0
		for scanner.Scan() {
			line := scanner.Text()
			lineIndent := len(yamlUtils.ExtractIndentationOfLine(line))
			if lineCounter < resourceLinesRange.Start+1 {
				lineCounter++
				continue
			}
			if lineCounter > resourceLinesRange.End || (tagsIndentSize > 0 && lineIndent <= tagsIndentSize) {
				tagsLines.End = lineCounter - 1
				break
			}
			if strings.TrimSpace(line) == FunctionTagsAttributeName+":" {
				tagsIndentSize = len(yamlUtils.ExtractIndentationOfLine(line))
				tagsLines.Start = lineCounter
				lineCounter++
				continue
			}
			lineCounter++
		}
	}
	if tagsLines.Start >= 0 && tagsLines.End == -1 {
		tagsLines.End = lineCounter - 1
	}
	return tagsLines
}
