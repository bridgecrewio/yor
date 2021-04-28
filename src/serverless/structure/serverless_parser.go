package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/tagging/tags"
	"bridgecrewio/yor/src/common/types"
	"bridgecrewio/yor/src/common/utils"
	"bufio"
	"fmt"
	"github.com/awslabs/goformation/v4/cloudformation"
	"io"
	"math"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	"io/ioutil"
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
		Functions    interface{}       `yaml:"functions"`
		Resources    struct {
			Resources *cloudformation.Template `yaml:"Resources"`
		} `yaml:"resources"`
	} `yaml:"provider"`
}

type ServerlessParser struct {
	YamlParser types.YamlParser
	Template   *ServerlessTemplate
}

func (p *ServerlessParser) Init(rootDir string, _ map[string]string) {
	p.YamlParser.RootDir = rootDir
	p.YamlParser.FileToResourcesLines = make(map[string]common.Lines)
	p.Template = &ServerlessTemplate{}
}

func (p *ServerlessParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *ServerlessParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension}
}

func (p *ServerlessParser) ParseFile(filePath string) ([]common.IBlock, error) {
	parsedBlocks := make([]common.IBlock, 0)
	template, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("There was an error processing the serverless template: %s", err))
	}
	err = yaml.Unmarshal(template, p.Template)
	if err != nil {
		logger.Error(fmt.Sprintf("Unmarshal: %s", err), "SILENT")
	}

	if err != nil || template == nil {
		logger.Error(fmt.Sprintf("There was an error processing the serverless template: %s", err), "SILENT")

	}
	//cfnStackTagsResource := p.template.Provider.CFNTags
	functions := p.Template.Provider.Functions
	functionsMap := functions.(map[interface{}]interface{})
	resourceNames := make([]string, 0)
	var resourceNamesToLines map[string]*common.Lines
	for funcName, _ := range functionsMap {
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
		var existingTags []tags.ITag
		funcRange := functionsMap[funcName].(map[interface{}]interface{})
		for key, val := range funcRange {
			lines := resourceNamesToLines[funcName]
			switch key {
			case FunctionTagsAttributeName:
				funcTags := val.(map[interface{}]interface{})
				for tagKey, tagValue := range funcTags {
					existingTags = append(existingTags, &tags.Tag{
						Key:   tagKey.(string),
						Value: tagValue.(string),
					})
				}
				tagsLines := p.extractLines(filePath, lines, resourceNames)
				rawBlock := funcRange
				minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
				maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))
				slsBlock := &ServerlessBlock{
					Block: common.Block{
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
				parsedBlocks = append(parsedBlocks, slsBlock)
			}
			p.YamlParser.FileToResourcesLines[filePath] = common.Lines{Start: minResourceLine, End: maxResourceLine}
		}
	}
	return parsedBlocks, nil
}

func (p *ServerlessParser) extractLines(filePath string, lines *common.Lines, resourceNames []string) common.Lines {
	tagsLines := p.getTagsLines(filePath, lines, resourceNames)
	return tagsLines
}

func (p *ServerlessParser) WriteFile(readFilePath string, blocks []common.IBlock, writeFilePath string) error {
	updatedBlocks := utils.EncodeBlocksToYaml(readFilePath, blocks, writeFilePath, FunctionTagsAttributeName, p.YamlParser.FileToResourcesLines[readFilePath])
	return utils.WriteYAMLFile(readFilePath, updatedBlocks, writeFilePath, p.YamlParser.FileToResourcesLines[readFilePath], FunctionTagsAttributeName)
}

func MapResourcesLineYAML(filePath string, resourceNames []string) map[string]*common.Lines {
	resourceToLines := make(map[string]*common.Lines)
	computedResources := make(map[string]bool, 0)
	for _, resourceName := range resourceNames {
		// initialize a map between resource name and its lines in file
		resourceToLines[resourceName] = &common.Lines{Start: -1, End: -1}
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

		if !readFunctions && latestResourceName != "" {
			// if we already read all the resources set the end line of the last resource and stop iterating the file
			resourceToLines[latestResourceName].End = lineCounter - 1
			computedResources[latestResourceName] = true
			break
		}
	}
	latestResourceName = ""
	funcLineIndentation := -1
	scanner = bufio.NewScanner(file)
	file.Seek(0, io.SeekStart)
	lineCounter = 0
	doneFunctions := false
	for scanner.Scan() {
		if !doneFunctions || len(computedResources) < len(resourceNames) {
			lineCounter++
			line := scanner.Text()
			sanitizedLine := strings.ReplaceAll(strings.TrimSpace(line), ":", "")
			lineIndentation := len(utils.ExtractIndentationOfLine(line))
			for _, resourceName := range resourceNames {
				if readFunctions {
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
						case int(funcLineIndentation):
							if resourceToLines[latestResourceName].End == -1 || utils.InSlice(resourceNames, sanitizedLine) {
								resourceToLines[latestResourceName].End = lineCounter - 1
								if sanitizedLine != latestResourceName {
									computedResources[latestResourceName] = true
									latestResourceName = sanitizedLine
								}
							}
							break
						case int(functionsSectionlineIndentation):
							// End functions sections
							resourceToLines[latestResourceName].End = lineCounter
							computedResources[latestResourceName] = true
							break
						default:
							if lineIndentation <= int(functionsSectionlineIndentation) && line != "" {
								resourceToLines[latestResourceName].End = lineCounter - 1
								computedResources[latestResourceName] = true
								doneFunctions = true
								break
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

func (p *ServerlessParser) getTagsLines(filePath string, resourceLinesRange *common.Lines, resourceNames []string) common.Lines {
	nonFoundLines := common.Lines{Start: -1, End: -1}
	switch utils.GetFileFormat(filePath) {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		//#nosec G304
		file, err := os.Open(filePath)
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
			return nonFoundLines
		}
		scanner := bufio.NewScanner(file)
		defer func() {
			_ = file.Close()
		}()
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
		linesInResource := utils.FindTagsLinesYAML(resourceLinesText, FunctionTagsAttributeName)
		numTags := linesInResource.End - linesInResource.Start
		return common.Lines{Start: linesInResource.Start + resourceLinesRange.Start, End: resourceLinesRange.End - numTags + 1}
	default:
		return common.Lines{Start: -1, End: -1}
	}
}
