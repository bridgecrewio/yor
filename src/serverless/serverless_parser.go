package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"bufio"
	"fmt"
	"github.com/awslabs/goformation/v4/cloudformation"
	"os"
	"strings"

	goformationtags "github.com/awslabs/goformation/v4/cloudformation/tags"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"reflect"
)

const ProviderTagsAttributeName = "tags"
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
	rootDir              string
	fileToResourcesLines map[string]common.Lines
	template             *ServerlessTemplate
}

func (p *ServerlessParser) Init(rootDir string, _ map[string]string) {
	p.rootDir = rootDir
	p.fileToResourcesLines = make(map[string]common.Lines)
	p.template = &ServerlessTemplate{}
}

func (p *ServerlessParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *ServerlessParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension}
}

func (p *ServerlessParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	template, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("There was an error processing the serverless template: %s", err))
	}
	err = yaml.Unmarshal([]byte(template), p.template)
	if err != nil {
		logger.Error(fmt.Sprintf("Unmarshal: %s", err), "SILENT")
	}

	if err != nil || template == nil {
		logger.Error(fmt.Sprintf("There was an error processing the serverless template: %s", err), "SILENT")

	}
	cfnStackTagsResource := p.template.Provider.CFNTags
	functions := p.template.Provider.Functions
	fmt.Println(functions, cfnStackTagsResource)
	value := reflect.ValueOf(functions)
	resourceNames := make([]string, 0)

	if value.Kind() == reflect.Map {
		for _, funcNameRef := range value.MapKeys() {
			funcName := funcNameRef.Elem().String()
			resourceNames = append(resourceNames, funcName)
			funcRange := value.MapIndex(funcNameRef).Elem().MapKeys()
			for _, keyRef := range funcRange {
				key := keyRef.Elem().String()
				val := value.MapIndex(funcNameRef).Elem().MapKeys()
				fmt.Println(funcName, key, val)
				switch key {
				case "tags":
					tagsRange := value.MapIndex(funcNameRef).Elem().MapIndex(keyRef).Elem()
					for _, tagKeyRef := range tagsRange.MapKeys() {
						tagKey := tagKeyRef.Elem().String()
						tagValue := tagsRange.MapIndex(tagKeyRef).Elem().String()
						fmt.Println(tagKey, tagValue)
					}
				}
			}
		}
	}
	var resourceNamesToLines map[string]*common.Lines
	switch common.GetFileFormat(filePath) {
	case common.YmlFileType.FileFormat, common.YamlFileType.FileFormat:
		resourceNamesToLines = MapResourcesLineYAML(filePath, resourceNames)
	default:
		return nil, fmt.Errorf("unsupported file type %s", common.GetFileFormat(filePath))
	}
	fmt.Println(resourceNamesToLines)

	return nil, nil
	// TODO
	//resourceNames := make([]string, 0)
	//if template != nil {
	//	for resourceName := range template.Resources {
	//		resourceNames = append(resourceNames, resourceName)
	//	}
	//
	//	var resourceNamesToLines map[string]*common.Lines
	//	switch common.GetFileFormat(filePath) {
	//	case common.YmlFileType.FileFormat, common.YamlFileType.FileFormat:
	//		resourceNamesToLines = MapResourcesLineYAML(filePath, resourceNames)
	//	default:
	//		return nil, fmt.Errorf("unsupported file type %s", common.GetFileFormat(filePath))
	//	}
	//
	//	minResourceLine := math.MaxInt8
	//	maxResourceLine := 0
	//	parsedBlocks := make([]structure.IBlock, 0)
	//	for resourceName := range template.Resources {
	//		resource := template.Resources[resourceName]
	//		lines := resourceNamesToLines[resourceName]
	//		isTaggable, tagsValue := common.StructContainsProperty(resource, TagsAttributeName)
	//		tagsLines := common.Lines{Start: -1, End: -1}
	//		var existingTags []tags.ITag
	//		if isTaggable {
	//			tagsLines, existingTags = p.extractTagsAndLines(filePath, lines, tagsValue)
	//		}
	//		minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
	//		maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))
	//
	//		cfnBlock := &ServerlessBlock{
	//			Block: structure.Block{
	//				FilePath:          filePath,
	//				ExitingTags:       existingTags,
	//				RawBlock:          resource,
	//				IsTaggable:        isTaggable,
	//				TagsAttributeName: TagsAttributeName,
	//			},
	//			lines:    *lines,
	//			name:     resourceName,
	//			tagLines: tagsLines,
	//		}
	//		parsedBlocks = append(parsedBlocks, cfnBlock)
	//	}
	//
	//	p.fileToResourcesLines[filePath] = common.Lines{Start: minResourceLine, End: maxResourceLine}
	//
	//	return parsedBlocks, nil
	//}
	return nil, err
}

func (p *ServerlessParser) extractTagsAndLines(filePath string, lines *common.Lines, tagsValue reflect.Value) (common.Lines, []tags.ITag) {
	tagsLines := p.getTagsLines(filePath, lines)
	existingTags := p.GetExistingTags(tagsValue)
	return tagsLines, existingTags
}

func (p *ServerlessParser) GetExistingTags(tagsValue reflect.Value) []tags.ITag {
	existingTags := make([]goformationtags.Tag, 0)
	if tagsValue.Kind() == reflect.Slice {
		existingTags = tagsValue.Interface().([]goformationtags.Tag)
	}

	iTags := make([]tags.ITag, 0)
	for _, goformationTag := range existingTags {
		tag := &tags.Tag{Key: goformationTag.Key, Value: goformationTag.Value}
		iTags = append(iTags, tag)
	}

	return iTags
}

func (p *ServerlessParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	fileFormat := common.GetFileFormat(readFilePath)
	switch fileFormat {
	case common.YamlFileType.Extension, common.YmlFileType.Extension:
		for _, block := range blocks {
			serverlessBlock, ok := block.(*ServerlessBlock)
			if !ok {
				logger.Warning("failed to convert block to ServerlessBlock")
				continue
			}
			serverlessBlock.UpdateTags()
		}
		return structure.WriteYAMLFile(readFilePath, blocks, writeFilePath, p.fileToResourcesLines[readFilePath], ProviderTagsAttributeName)

	default:
		logger.Warning(fmt.Sprintf("unsupported file type %s", fileFormat))
		return nil
	}
}

func MapResourcesLineYAML(filePath string, resourceNames []string) map[string]*common.Lines {
	resourceToLines := make(map[string]*common.Lines)
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

	readResources := false
	lineCounter := 0
	latestResourceName := ""
	// iterate file line by line
	for scanner.Scan() {
		lineCounter++
		line := scanner.Text()

		// make sure we look for resources names only under the Resources section
		foundSectionIndex := -1
		for i, templateSectionName := range templateSections {
			if strings.Contains(line, templateSectionName) {
				foundSectionIndex = i
				readResources = templateSectionName == "Resources"
				break
			}
		}
		if foundSectionIndex >= 0 {
			// if this is a section line, check if we're done reading resources, otherwise remove the section name

			if !readResources && latestResourceName != "" {
				// if we already read all the resources set the end line of the last resource and stop iterating the file
				resourceToLines[latestResourceName].End = lineCounter - 1
				break
			}
			// remove found section to avoid searching for it once it was found
			templateSections = append(templateSections[:foundSectionIndex], templateSections[foundSectionIndex+1:]...)
			continue
		}

		if readResources {
			foundResourceIndex := -1
			for i, resourceName := range resourceNames {
				if strings.Contains(line, resourceName) {
					if latestResourceName != "" {
						// set the end line of the previous resource
						resourceToLines[latestResourceName].End = lineCounter - 1
					}

					foundResourceIndex = i
					resourceToLines[resourceName].Start = lineCounter
					latestResourceName = resourceName
					break
				}
			}
			if foundResourceIndex >= 0 {
				// remove found resource name to avoid searching for it once it was found
				resourceNames = append(resourceNames[:foundResourceIndex], resourceNames[foundResourceIndex+1:]...)
				continue
			}
		}
	}
	if latestResourceName != "" && resourceToLines[latestResourceName].End == -1 {
		// in case we reached the end of the file without setting the end line of the last resource
		resourceToLines[latestResourceName].End = lineCounter
	}

	return resourceToLines
}

func (p *ServerlessParser) getTagsLines(filePath string, resourceLinesRange *common.Lines) common.Lines {
	nonFoundLines := common.Lines{Start: -1, End: -1}
	switch common.GetFileFormat(filePath) {
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
		for scanner.Scan() {
			if lineCounter > resourceLinesRange.End {
				break
			}
			if lineCounter >= resourceLinesRange.Start && lineCounter <= resourceLinesRange.End {
				resourceLinesText = append(resourceLinesText, scanner.Text())
			}
			lineCounter++
		}
		linesInResource := structure.FindTagsLinesYAML(resourceLinesText, ProviderTagsAttributeName)
		return common.Lines{Start: linesInResource.Start + resourceLinesRange.Start, End: linesInResource.End + resourceLinesRange.End}
	default:
		return common.Lines{Start: -1, End: -1}
	}
}
