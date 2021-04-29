package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/tagging/tags"
	"bridgecrewio/yor/src/common/types"
	"bridgecrewio/yor/src/common/utils"
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/awslabs/goformation/v4"
	goformationTags "github.com/awslabs/goformation/v4/cloudformation/tags"

	"reflect"
)

type CloudformationParser struct {
	*types.YamlParser
}

func (p *CloudformationParser) StructContainsProperty(s interface{}, property string) (bool, reflect.Value) {
	var field reflect.Value
	sValue := reflect.ValueOf(s)

	// Check if the passed interface is a pointer
	if sValue.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		field = sValue.FieldByName(property)
	} else {
		// 'dereference' with Elem() and get the field by name
		field = sValue.Elem().FieldByName(property)
	}

	if !field.IsValid() {
		return false, field
	}

	return true, field
}

const TagsAttributeName = "Tags"

var TemplateSections = []string{"AWSTemplateFormatVersion", "Transform", "Description", "Metadata", "Parameters", "Mappings", "Conditions", "Outputs", "Resources"}

func (p *CloudformationParser) Init(rootDir string, _ map[string]string) {
	p.YamlParser = &types.YamlParser{
		RootDir:              rootDir,
		FileToResourcesLines: make(map[string]common.Lines),
	}
}

func (p *CloudformationParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *CloudformationParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension, common.CFTFileType.Extension}
}

func (p *CloudformationParser) ParseFile(filePath string) ([]common.IBlock, error) {
	template, err := goformation.Open(filePath)
	if err != nil || template == nil {
		logger.Warning(fmt.Sprintf("There was an error processing the cloudformation template: %s", err))
	}

	resourceNames := make([]string, 0)
	if template != nil {
		for resourceName := range template.Resources {
			resourceNames = append(resourceNames, resourceName)
		}

		var resourceNamesToLines map[string]*common.Lines
		switch utils.GetFileFormat(filePath) {
		case common.YmlFileType.FileFormat, common.YamlFileType.FileFormat:
			resourceNamesToLines = MapResourcesLineYAML(filePath, resourceNames)
		default:
			return nil, fmt.Errorf("unsupported file type %s", utils.GetFileFormat(filePath))
		}

		minResourceLine := math.MaxInt8
		maxResourceLine := 0
		parsedBlocks := make([]common.IBlock, 0)
		for resourceName := range template.Resources {
			resource := template.Resources[resourceName]
			lines := resourceNamesToLines[resourceName]
			isTaggable, tagsValue := p.StructContainsProperty(resource, TagsAttributeName)
			tagsLines := common.Lines{Start: -1, End: -1}
			var existingTags []tags.ITag
			if isTaggable {
				tagsLines, existingTags = p.extractTagsAndLines(filePath, lines, tagsValue)
			}
			minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
			maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))

			cfnBlock := &CloudformationBlock{
				Block: common.Block{
					FilePath:          filePath,
					ExitingTags:       existingTags,
					RawBlock:          resource,
					IsTaggable:        isTaggable,
					TagsAttributeName: TagsAttributeName,
					Lines:             *lines,
					TagLines:          tagsLines,
				},
				Name: resourceName,
			}
			parsedBlocks = append(parsedBlocks, cfnBlock)
		}

		p.FileToResourcesLines[filePath] = common.Lines{Start: minResourceLine, End: maxResourceLine}

		return parsedBlocks, nil
	}
	return nil, err
}

func (p *CloudformationParser) extractTagsAndLines(filePath string, lines *common.Lines, tagsValue reflect.Value) (common.Lines, []tags.ITag) {
	tagsLines := p.getTagsLines(filePath, lines)
	existingTags := p.GetExistingTags(tagsValue)
	return tagsLines, existingTags
}

func (p *CloudformationParser) GetExistingTags(tagsValue reflect.Value) []tags.ITag {
	existingTags := make([]goformationTags.Tag, 0)
	if tagsValue.Kind() == reflect.Slice {
		existingTags = tagsValue.Interface().([]goformationTags.Tag)
	}

	iTags := make([]tags.ITag, 0)
	for _, goformationTag := range existingTags {
		tag := &tags.Tag{Key: goformationTag.Key, Value: goformationTag.Value}
		iTags = append(iTags, tag)
	}

	return iTags
}

func (p *CloudformationParser) WriteFile(readFilePath string, blocks []common.IBlock, writeFilePath string) error {
	err := utils.EncodeBlocksToYaml(readFilePath, blocks, writeFilePath, TagsAttributeName, p.YamlParser.FileToResourcesLines[readFilePath])
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to encode %s", readFilePath))
	}
	return utils.WriteYAMLFile(readFilePath, blocks, writeFilePath, p.FileToResourcesLines[readFilePath], TagsAttributeName)
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

func (p *CloudformationParser) getTagsLines(filePath string, resourceLinesRange *common.Lines) common.Lines {
	nonFoundLines := common.Lines{Start: -1, End: -1}
	switch utils.GetFileFormat(filePath) {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		scanner, _ := utils.GetFileScanner(filePath, &nonFoundLines)
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
		linesInResource := utils.FindTagsLinesYAML(resourceLinesText, TagsAttributeName)
		return common.Lines{Start: linesInResource.Start + resourceLinesRange.Start, End: linesInResource.End + resourceLinesRange.End}
	default:
		return common.Lines{Start: -1, End: -1}
	}
}
