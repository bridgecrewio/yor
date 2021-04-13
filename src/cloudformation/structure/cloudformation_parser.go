package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/awslabs/goformation/v4"
	goformation_tags "github.com/awslabs/goformation/v4/cloudformation/tags"

	"reflect"
)

const TagsAttributeName = "Tags"

var TemplateSections = []string{"AWSTemplateFormatVersion", "Transform", "Description", "Metadata", "Parameters", "Mappings", "Conditions", "Outputs", "Resources"}

type CloudformationParser struct {
	rootDir              string
	fileToResourcesLines map[string]common.Lines
}

func (p *CloudformationParser) Init(rootDir string, _ map[string]string) {
	p.rootDir = rootDir
	p.fileToResourcesLines = make(map[string]common.Lines)
}

func (p *CloudformationParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *CloudformationParser) GetAllowedFileTypes() []string {
	return []string{".yaml", "yml"}
}

func (p *CloudformationParser) ParseFile(filePath string) ([]structure.IBlock, error) {
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
		switch common.GetFileFormat(filePath) {
		case "yaml", "yml":
			resourceNamesToLines = MapResourcesLineYAML(filePath, resourceNames)
		default:
			return nil, fmt.Errorf("unsupported file type %s", common.GetFileFormat(filePath))
		}

		minResourceLine := math.MaxInt8
		maxResourceLine := 0
		parsedBlocks := make([]structure.IBlock, 0)
		for resourceName := range template.Resources {
			resource := template.Resources[resourceName]
			lines := resourceNamesToLines[resourceName]
			isTaggable, tagsValue := common.StructContainsProperty(resource, TagsAttributeName)
			var tagsLines common.Lines
			var existingTags []tags.ITag
			if isTaggable {
				tagsLines, existingTags = p.extractTagsAndLines(filePath, lines, tagsValue)
			}
			minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
			maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))

			cfnBlock := &CloudformationBlock{
				Block: structure.Block{
					FilePath:          filePath,
					ExitingTags:       existingTags,
					RawBlock:          resource,
					IsTaggable:        isTaggable,
					TagsAttributeName: TagsAttributeName,
				},
				lines:    *lines,
				name:     resourceName,
				tagLines: tagsLines,
			}
			parsedBlocks = append(parsedBlocks, cfnBlock)
		}

		p.fileToResourcesLines[filePath] = common.Lines{Start: minResourceLine, End: maxResourceLine}

		return parsedBlocks, nil
	}
	return nil, err
}

func (p *CloudformationParser) extractTagsAndLines(filePath string, lines *common.Lines, tagsValue reflect.Value) (common.Lines, []tags.ITag) {
	var existingTags []tags.ITag
	tagsLines := common.Lines{Start: -1, End: -1}
	tagsLines = p.getTagsLines(filePath, lines)
	existingTags = p.GetExistingTags(tagsValue)
	return tagsLines, existingTags
}

func (p *CloudformationParser) GetExistingTags(tagsValue reflect.Value) []tags.ITag {
	existingTags := make([]goformation_tags.Tag, 0)
	if tagsValue.Kind() == reflect.Slice {
		existingTags = tagsValue.Interface().([]goformation_tags.Tag)
	}

	iTags := make([]tags.ITag, 0)
	for _, goformationTag := range existingTags {
		tag := &tags.Tag{Key: goformationTag.Key, Value: goformationTag.Value}
		iTags = append(iTags, tag)
	}

	return iTags
}

func (p *CloudformationParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	fileFormat := common.GetFileFormat(readFilePath)
	switch fileFormat {
	case "yaml":
		for _, block := range blocks {
			cloudformationBlock, ok := block.(*CloudformationBlock)
			if !ok {
				logger.Warning("failed to convert block to CloudformationBlock")
				continue
			}
			cloudformationBlock.UpdateTags()
		}
		return structure.WriteYAMLFile(readFilePath, blocks, writeFilePath, p.fileToResourcesLines[readFilePath], TagsAttributeName)
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

func (p *CloudformationParser) getTagsLines(filePath string, resourceLinesRange *common.Lines) common.Lines {
	nonFoundLines := common.Lines{Start: -1, End: -1}
	switch common.GetFileFormat(filePath) {
	case "yaml":
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
		linesInResource := structure.FindTagsLinesYAML(resourceLinesText, TagsAttributeName)
		return common.Lines{Start: linesInResource.Start + resourceLinesRange.Start, End: linesInResource.End + resourceLinesRange.End}
	default:
		return common.Lines{Start: -1, End: -1}
	}
}
