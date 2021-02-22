package structure

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"bufio"
	"fmt"
	"github.com/awslabs/goformation/v4"
	goformation_tags "github.com/awslabs/goformation/v4/cloudformation/tags"

	"os"
	"reflect"
	"strings"
)

const TagsAttributeName = "Tags"

var TemplateSections = []string{"AWSTemplateFormatVersion", "Transform", "Description", "Metadata", "Parameters", "Mappings", "Conditions", "Outputs", "Resources"}

type CloudformationParser struct {
	rootDir string
}

func (p *CloudformationParser) Init(rootDir string, _ map[string]string) {
	p.rootDir = rootDir
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

	parsedBlocks := make([]structure.IBlock, 0)
	for resourceName := range template.Resources {
		resource := template.Resources[resourceName]
		isTaggable, tagsValue := common.StructContainsProperty(resource, TagsAttributeName)
		var existingTags []tags.ITag
		if isTaggable {
			existingTags = p.GetExistingTags(tagsValue)
		}
		cfnBlock := &CloudformationBlock{
			Block: structure.Block{
				FilePath:          filePath,
				ExitingTags:       existingTags,
				RawBlock:          resource,
				IsTaggable:        isTaggable,
				TagsAttributeName: TagsAttributeName,
			},
			lines: resourceNamesToLines[resourceName],
			name:  resourceName,
		}
		parsedBlocks = append(parsedBlocks, cfnBlock)
	}

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

func mapResourcesLineYAML(filePath string, resourceNames []string) map[string][]int {
	resourceToLines := make(map[string][]int)
	for _, resourceName := range resourceNames {
		// initialize a map between resource name and its lines in file
		resourceToLines[resourceName] = []int{-1, -1}
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
				resourceToLines[latestResourceName][1] = lineCounter - 1
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
						resourceToLines[latestResourceName][1] = lineCounter - 1
					}

					foundResourceIndex = i
					resourceToLines[resourceName][0] = lineCounter
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
	if latestResourceName != "" && resourceToLines[latestResourceName][1] == -1 {
		// in case we reached the end of the file without setting the end line of the last resource
		resourceToLines[latestResourceName][1] = lineCounter
	}

	return resourceToLines
}

func (p *CloudformationParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {

	return nil
}
