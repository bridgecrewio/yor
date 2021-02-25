package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"bufio"
	"fmt"
	"github.com/awslabs/goformation/v4"
	"github.com/awslabs/goformation/v4/cloudformation"
	goformation_tags "github.com/awslabs/goformation/v4/cloudformation/tags"
	"github.com/sanathkr/yaml"
	"io/ioutil"
	"math"

	"os"
	"reflect"
	"strings"
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
	return []string{".yaml"}
}

func (p *CloudformationParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	template, err := goformation.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("There was an error processing the cloudformation template: %s", err))
	}

	resourceNames := make([]string, 0)
	for resourceName := range template.Resources {
		resourceNames = append(resourceNames, resourceName)
	}
	resourceNamesToLines := mapResourcesLineYAML(filePath, resourceNames)
	minResourceLine := math.MaxInt8
	maxResourceLine := 0
	parsedBlocks := make([]structure.IBlock, 0)
	for resourceName := range template.Resources {
		resource := template.Resources[resourceName]
		isTaggable, tagsValue := common.StructContainsProperty(resource, TagsAttributeName)
		var existingTags []tags.ITag
		if isTaggable {
			existingTags = p.GetExistingTags(tagsValue)
		}
		lines := resourceNamesToLines[resourceName]
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
			lines: *lines,
			name:  resourceName,
		}
		parsedBlocks = append(parsedBlocks, cfnBlock)
	}

	p.fileToResourcesLines[filePath] = common.Lines{Start: minResourceLine, End: maxResourceLine}

	return parsedBlocks, nil
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

func mapResourcesLineYAML(filePath string, resourceNames []string) map[string]*common.Lines {
	resourceToLines := make(map[string]*common.Lines)
	for _, resourceName := range resourceNames {
		// initialize a map between resource name and its lines in file
		resourceToLines[resourceName] = &common.Lines{Start: -1, End: -1}
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
		return nil
	}
	scanner := bufio.NewScanner(file)
	defer file.Close()

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

func (p *CloudformationParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	// read file bytes
	fileFormat := common.GetFileFormat(readFilePath)
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}

	originLines := common.GetLinesFromBytes(originFileSrc)

	resourcesStart := p.fileToResourcesLines[readFilePath].Start
	resourcesEnd := p.fileToResourcesLines[readFilePath].End
	resourcesIndent := extractIndentationOfLine(originLines[resourcesStart])
	linesBeforeResources := originLines[:resourcesStart-1]
	linesAfterResources := originLines[resourcesEnd:]
	resourcesLines := make([]string, 0)
	for _, resourceBlock := range blocks {
		cloudformationBlock, ok := resourceBlock.(*CloudformationBlock)
		if !ok {
			logger.Warning("failed to convert block to CloudformationBlock")
			continue
		}
		cloudformationBlock.UpdateTags()
		oldBlockLines := resourceBlock.GetLines()
		oldResourceLines := originLines[oldBlockLines.Start-1 : oldBlockLines.End]
		oldResourceTagLines := findTagsLinesYAML(oldResourceLines)
		tagsIndent := extractIndentationOfLine(oldResourceLines[oldResourceTagLines.Start+1])

		newResourceLines := p.GetResourcesLines("", fileFormat, resourceBlock)

		newResourceTagLines := findTagsLinesYAML(newResourceLines)
		resourcesLines = append(resourcesLines, oldResourceLines[:oldResourceTagLines.Start]...)
		resourcesLines = append(resourcesLines, resourcesIndent+newResourceLines[newResourceTagLines.Start])
		resourcesLines = append(resourcesLines, indentLines(newResourceLines[newResourceTagLines.Start+1:newResourceTagLines.End+1], tagsIndent)...)
		resourcesLines = append(resourcesLines, oldResourceLines[oldResourceTagLines.End+1:]...)
	}

	allLines := append(linesBeforeResources, resourcesLines...)
	allLines = append(allLines, linesAfterResources...)
	linesText := strings.Join(allLines, "\n")

	err = ioutil.WriteFile(writeFilePath, []byte(linesText), 0644)

	return err
}

func extractIndentationOfLine(textLine string) string {
	indent := ""
	for _, c := range textLine {
		if c != ' ' {
			break
		}
		indent += " "
	}

	return indent
}

func (p *CloudformationParser) GetResourcesLines(indent string, fileExtension string, block structure.IBlock) []string {
	cfnResource := block.GetRawBlock().(cloudformation.Resource)
	switch fileExtension {
	case "yaml":
		return p.getYAMLLines(indent, cfnResource)
	default:
		logger.Warning(fmt.Sprintf("unsupported file type %s", fileExtension))
	}

	return nil
}

func (p *CloudformationParser) getYAMLLines(indent string, cfnResource cloudformation.Resource) []string {
	yamlBytes, err := yaml.Marshal(cfnResource)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal cloudformation resource to yaml: %s", err))
	}
	textLines := common.GetLinesFromBytes(yamlBytes)
	indentedLines := make([]string, len(textLines))
	for i, originLine := range textLines {
		indentedLines[i] = indent + originLine
	}

	return indentedLines
}

func indentLines(textLines []string, indent string) []string {
	originIndent := extractIndentationOfLine(textLines[0])
	for i, originLine := range textLines {
		noLeadingWhitespace := originLine[len(originIndent):]
		textLines[i] = indent + noLeadingWhitespace
	}

	return textLines
}

func findTagsLinesYAML(textLines []string) common.Lines {
	tagsLines := common.Lines{Start: -1, End: len(textLines) - 1}
	indent := ""
	for i, line := range textLines {
		if strings.Contains(line, TagsAttributeName+":") {
			tagsLines.Start = i
			indent = extractIndentationOfLine(line)
		} else if extractIndentationOfLine(line) <= indent && !strings.HasPrefix(line, indent+"-") && tagsLines.Start >= 0 {
			tagsLines.End = i - 1
			return tagsLines
		}
	}

	return tagsLines
}
