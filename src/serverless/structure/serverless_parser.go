package structure

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/types"
	"github.com/bridgecrewio/yor/src/common/utils"
	yamlUtils "github.com/bridgecrewio/yor/src/common/yaml"
	"gopkg.in/yaml.v2"
)

const FunctionTagsAttributeName = "tags"
const StackTagsAttributeName = "stackTags"

var TemplateSections = []string{"service", "provider", "tags", "stackTags", "resources", "functions", "region", "Resources"}

type ServerlessTemplate struct {
	Service  string `yaml:"service"`
	Provider struct {
		Name         string            `yaml:"name"`
		Runtime      string            `yaml:"runtime"`
		Region       string            `yaml:"region"`
		ProviderTags map[string]string `yaml:"tags"`
		CFNTags      map[string]string `yaml:"stackTags"`
	} `yaml:"provider"`
	Functions interface{} `yaml:"functions"`
	Resources struct {
		Resources *cloudformation.Template `yaml:"Resources"`
	} `yaml:"resources"`
}

type ServerlessParser struct {
	YamlParser types.YamlParser
	Template   *ServerlessTemplate
}

func (p *ServerlessParser) Init(rootDir string, _ map[string]string) {
	p.YamlParser.RootDir = rootDir
	p.YamlParser.FileToResourcesLines = make(map[string]structure.Lines)
	p.Template = &ServerlessTemplate{}
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
	template, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("There was an error processing the serverless template: %s", err))
	}
	err = yaml.Unmarshal(template, p.Template)
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
	functions := p.Template.Functions
	functionsMap := functions.(map[interface{}]interface{})
	resourceNames := make([]string, 0)
	var resourceNamesToLines map[string]*structure.Lines
	for funcName := range functionsMap {
		resourceNames = append(resourceNames, funcName.(string))
	}
	switch utils.GetFileFormat(filePath) {
	case common.YmlFileType.FileFormat, common.YamlFileType.FileFormat:
		resourceNamesToLines = MapResourcesLineYAML(filePath, resourceNames)
	default:
		return nil, fmt.Errorf("unsupported file type %s", utils.GetFileFormat(filePath))
	}
	minResourceLine := math.MaxInt8
	maxResourceLine := 0
	for _, funcName := range resourceNames {
		var rawBlock interface{}
		var slsBlock *ServerlessBlock
		tagsExist := false
		var existingTags = make([]tags.ITag, 0)
		tagsLines := structure.Lines{Start: -1, End: -1}
		funcRange := functionsMap[funcName].(map[interface{}]interface{})
		var lines *structure.Lines
		rawBlock = funcRange
		for key, val := range funcRange {
			lines = resourceNamesToLines[funcName]
			minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
			maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))
			if key == FunctionTagsAttributeName {
				tagsLines = p.getTagsLines(filePath, lines, resourceNames)
				tagsExist = true
				funcTags := val.(map[interface{}]interface{})
				for tagKey, tagValue := range funcTags {
					existingTags = append(existingTags, &tags.Tag{
						Key:   tagKey.(string),
						Value: tagValue.(string),
					})
				}
			}
		}
		if !tagsExist {
			rawBlock.(map[interface{}]interface{})[FunctionTagsAttributeName] = make([]map[string]string, 0)
			slsBlock = &ServerlessBlock{
				Block: structure.Block{
					FilePath:          filePath,
					ExitingTags:       existingTags,
					RawBlock:          rawBlock,
					IsTaggable:        true,
					TagsAttributeName: FunctionTagsAttributeName,
					Lines:             *lines,
					TagLines:          tagsLines,
				},
				Name: funcName,
			}
		} else {
			slsBlock = &ServerlessBlock{
				Block: structure.Block{
					FilePath:          filePath,
					ExitingTags:       existingTags,
					RawBlock:          rawBlock,
					IsTaggable:        true,
					TagsAttributeName: FunctionTagsAttributeName,
					Lines:             *lines,
					TagLines:          tagsLines,
				},
				Name: funcName,
			}
		}

		parsedBlocks = append(parsedBlocks, slsBlock)
		p.YamlParser.FileToResourcesLines[filePath] = structure.Lines{Start: minResourceLine, End: maxResourceLine}

	}
	return parsedBlocks, nil
}

func (p *ServerlessParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	updatedBlocks := yamlUtils.EncodeBlocksToYaml(readFilePath, blocks)
	return yamlUtils.WriteYAMLFile(readFilePath, updatedBlocks, writeFilePath, p.YamlParser.FileToResourcesLines[readFilePath], FunctionTagsAttributeName)
}

func MapResourcesLineYAML(filePath string, resourceNames []string) map[string]*structure.Lines {
	resourceToLines := make(map[string]*structure.Lines)
	computedResources := make(map[string]bool)
	for _, resourceName := range resourceNames {
		// initialize a map between resource name and its lines in file
		resourceToLines[resourceName] = &structure.Lines{Start: -1, End: -1}
	}
	// #nosec G304
	file, err := os.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
		return nil
	}
	scanner := bufio.NewScanner(file)
	defer func() {
		_ = file.Close()
	}()

	// deep copy TemplateSections to allow modifying it safely
	templateSections := make([]string, len(TemplateSections))
	copy(templateSections, TemplateSections)

	readFunctions := false
	functionsSectionlineIndentation := -1
	lineCounter := 0
	latestResourceName := ""
	// iterate file line by line
	for scanner.Scan() {
		lineCounter++
		line := scanner.Text()

		for _, templateSectionName := range templateSections {
			if strings.Contains(line, templateSectionName) {
				if !readFunctions {
					readFunctions = templateSectionName == "functions"
				}
				if readFunctions {
					functionsSectionlineIndentation = len(utils.ExtractIndentationOfLine(line))
				}
			}
		}
		if readFunctions {
			for _, resourceName := range resourceNames {
				if strings.Contains(line, resourceName) {
					latestResourceName = resourceName
				}
			}
		}

		if readFunctions && latestResourceName != "" {
			// if we already read all the resources set the end line of the last resource and stop iterating the file
			resourceToLines[latestResourceName].End = lineCounter - 1
			computedResources[latestResourceName] = true
			break
		}
	}
	latestResourceName = ""
	funcLineIndentation := -1
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		logger.Error(err.Error())
	}
	scanner = bufio.NewScanner(file)
	lineCounter = 0
	readFunctions = false
	for scanner.Scan() {
		if !readFunctions || len(computedResources) < len(resourceNames) {
			lineCounter++
			line := scanner.Text()
			sanitizedLine := strings.ReplaceAll(strings.TrimSpace(line), ":", "")
			lineIndentation := len(utils.ExtractIndentationOfLine(line))
			for _, resourceName := range resourceNames {
				if !readFunctions {
					if sanitizedLine == resourceName {
						funcLineIndentation = lineIndentation
						if latestResourceName != "" {
							// set the end line of the previous resource
							resourceToLines[latestResourceName].End = lineCounter - 1
							computedResources[latestResourceName] = true
						}
						resourceToLines[resourceName].Start = lineCounter
						latestResourceName = resourceName
					}
					if latestResourceName != "" {
						switch lineIndentation {
						case funcLineIndentation:
							if resourceToLines[latestResourceName].End == -1 || utils.InSlice(resourceNames, sanitizedLine) {
								resourceToLines[latestResourceName].End = lineCounter - 1
								if sanitizedLine != latestResourceName {
									computedResources[latestResourceName] = true
									latestResourceName = sanitizedLine
								}
							}
						case functionsSectionlineIndentation:
							// End functions sections
							resourceToLines[latestResourceName].End = lineCounter
							computedResources[latestResourceName] = true
						default:
							if lineIndentation <= functionsSectionlineIndentation && line != "" {
								resourceToLines[latestResourceName].End = lineCounter - 1
								computedResources[latestResourceName] = true
								readFunctions = true
							}
						}
					}
				}
			}
			if latestResourceName != "" && resourceToLines[latestResourceName].End == -1 {
				// in case we reached the end of the file without setting the end line of the last resource
				resourceToLines[latestResourceName].End = lineCounter
				computedResources[latestResourceName] = true
			}
		}
	}
	return resourceToLines
}
func isLineFunctionDefinition(line string, resourceNames []string) bool {
	sanitizedLine := strings.ReplaceAll(strings.TrimSpace(line), ":", "")
	return utils.InSlice(resourceNames, sanitizedLine)
}

func (p *ServerlessParser) getTagsLines(filePath string, resourceLinesRange *structure.Lines, resourceNames []string) structure.Lines {
	nonFoundLines := structure.Lines{Start: -1, End: -1}
	fileFormat := utils.GetFileFormat(filePath)
	switch fileFormat {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		file, scanner, _ := utils.GetFileScanner(filePath, &nonFoundLines)
		resourceLinesText := make([]string, 0)
		// iterate file line by line
		lineCounter := 0
		funcIndentLevel := -1
		for scanner.Scan() {
			line := scanner.Text()
			lineIndent := len(utils.ExtractIndentationOfLine(line))
			if lineCounter == resourceLinesRange.Start-1 {
				funcIndentLevel = lineIndent
			}
			if lineCounter > resourceLinesRange.End {
				break
			}
			if isLineFunctionDefinition(line, resourceNames) && lineCounter > resourceLinesRange.Start-1 {
				break
			}
			if lineCounter >= resourceLinesRange.Start && lineCounter <= resourceLinesRange.End && (lineIndent > funcIndentLevel) {
				resourceLinesText = append(resourceLinesText, line)
			}
			lineCounter++
		}
		linesInResource, _ := yamlUtils.FindTagsLinesYAML(resourceLinesText, FunctionTagsAttributeName)
		numTags := linesInResource.End - linesInResource.Start
		defer func() {
			_ = file.Close()
		}()
		return structure.Lines{Start: linesInResource.Start + resourceLinesRange.Start, End: resourceLinesRange.End - numTags + 1}
	default:
		return structure.Lines{Start: -1, End: -1}
	}
}
