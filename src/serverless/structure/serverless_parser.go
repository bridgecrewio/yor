package structure

import (
	"fmt"
	"math"
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
	Template   *serverless.Template
}

func (p *ServerlessParser) Init(rootDir string, _ map[string]string) {
	p.YamlParser.RootDir = rootDir
	p.YamlParser.FileToResourcesLines = make(map[string]structure.Lines)
	p.Template = &serverless.Template{}
}

func (p *ServerlessParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *ServerlessParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension}
}

func (p *ServerlessParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	parsedBlocks := make([]structure.IBlock, 0)
	fileFormat := utils.GetFileFormat(filePath)
	fileName := filepath.Base(filePath)
	if !(fileName == fmt.Sprintf("serverless.%s", fileFormat)) {
		return nil, nil
	}
	// #nosec G304 - file is from user
	template, err := goserverless.Open(filePath)
	p.Template = template
	if err != nil {
		logger.Warning(fmt.Sprintf("There was an error processing the serverless template: %s", err))
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Unmarshal: %s", err), "SILENT")
	}
	if p.Template.Functions == nil && p.Template.Resources.Resources == nil {
		return parsedBlocks, nil
	}

	if err != nil || template == nil {
		logger.Error(fmt.Sprintf("There was an error processing the serverless template: %s", err), "SILENT")

	}
	// cfnStackTagsResource := p.template.Provider.CFNTags
	resourceNames := make([]string, 0)
	var resourceNamesToLines map[string]*structure.Lines
	for funcName := range p.Template.Functions {
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
	var existingTags []tags.ITag
	for _, funcName := range resourceNames {
		var slsBlock *ServerlessBlock
		tagsLines := structure.Lines{Start: -1, End: -1}
		var lines *structure.Lines
		slsFunction := p.Template.Functions[funcName]
		lines = resourceNamesToLines[funcName]
		minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
		maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))
		if slsFunction.Tags != nil {
			tagsLines = p.getTagsLines(filePath, lines)
			for tagKey, tagValue := range slsFunction.Tags {
				existingTags = append(existingTags, &tags.Tag{Key: tagKey, Value: tagValue.(string)})
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
		p.YamlParser.FileToResourcesLines[filePath] = structure.Lines{Start: minResourceLine, End: maxResourceLine}

	}
	return parsedBlocks, nil
}

func (p *ServerlessParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	updatedBlocks := yamlUtils.EncodeBlocksToYaml(readFilePath, blocks)
	return yamlUtils.WriteYAMLFile(readFilePath, updatedBlocks, writeFilePath, p.YamlParser.FileToResourcesLines[readFilePath],
		FunctionTagsAttributeName, FunctionsSectionName)
}

func (p *ServerlessParser) getTagsLines(filePath string, resourceLinesRange *structure.Lines) structure.Lines {
	nonFoundLines := structure.Lines{Start: -1, End: -1}
	fileFormat := utils.GetFileFormat(filePath)
	tagsLines := structure.Lines{Start: -1, End: -1}
	switch fileFormat {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		file, scanner, _ := utils.GetFileScanner(filePath, &nonFoundLines)
		defer func() {
			_ = file.Close()
		}()
		// iterate file line by line
		lineCounter := 0
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
	return tagsLines
}
