package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/sanathkr/yaml"
)

func MapResourcesLineYAML(filePath string, resourceNames []string) map[string]*common.Lines {
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

func WriteYAMLFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, resourcesLinesRange common.Lines) error {
	// read file bytes
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}
	originLines := common.GetLinesFromBytes(originFileSrc)

	resourcesIndent := common.ExtractIndentationOfLine(originLines[resourcesLinesRange.Start])

	resourcesLines := make([]string, 0)
	for _, resourceBlock := range blocks {
		cloudformationBlock, ok := resourceBlock.(*CloudformationBlock)
		if !ok {
			logger.Warning("failed to convert block to CloudformationBlock")
			continue
		}
		cloudformationBlock.UpdateTags()
		newResourceLines := getYAMLLines(resourceBlock.GetRawBlock().(cloudformation.Resource))
		newResourceTagLineRange := findTagsLinesYAML(newResourceLines)

		oldResourceLinesRange := resourceBlock.GetLines()
		oldResourceLines := originLines[oldResourceLinesRange.Start-1 : oldResourceLinesRange.End]

		// if the block is not taggable, write it and continue
		if !resourceBlock.IsBlockTaggable() {
			resourcesLines = append(resourcesLines, oldResourceLines...)
			continue
		}

		oldResourceTagLines := findTagsLinesYAML(oldResourceLines)
		// if the resource don't contain Tags entry - create it
		if oldResourceTagLines.Start == -1 || oldResourceTagLines.End == -1 {
			// get the indentation of the property under the resource name
			tagAttributeIndent := common.ExtractIndentationOfLine(oldResourceLines[1])
			resourcesLines = append(resourcesLines, oldResourceLines...)                      // add all the existing resource data first
			resourcesLines = append(resourcesLines, tagAttributeIndent+TagsAttributeName+":") // add the 'Tags:' line
			// add the tags with extra indentation below the 'Tags:' line
			resourcesLines = append(resourcesLines, common.IndentLines(newResourceLines[newResourceTagLineRange.Start+1:newResourceTagLineRange.End+1], tagAttributeIndent+"  ")...)
			continue
		}

		oldTagsIndent := common.ExtractIndentationOfLine(oldResourceLines[oldResourceTagLines.Start+1])

		resourcesLines = append(resourcesLines, oldResourceLines[:oldResourceTagLines.Start]...)                                                                       // add all the resource's line before the tags
		resourcesLines = append(resourcesLines, resourcesIndent+newResourceLines[newResourceTagLineRange.Start])                                                       // add the 'Tags:' line
		resourcesLines = append(resourcesLines, common.IndentLines(newResourceLines[newResourceTagLineRange.Start+1:newResourceTagLineRange.End+1], oldTagsIndent)...) // add tags
		resourcesLines = append(resourcesLines, oldResourceLines[oldResourceTagLines.End+1:]...)                                                                       // add rest of resource lines
	}

	allLines := append(originLines[:resourcesLinesRange.Start-1], resourcesLines...)
	allLines = append(allLines, originLines[resourcesLinesRange.End:]...)
	linesText := strings.Join(allLines, "\n")

	err = ioutil.WriteFile(writeFilePath, []byte(linesText), 0644)

	return err
}

func getYAMLLines(cfnResource cloudformation.Resource) []string {
	yamlBytes, err := yaml.Marshal(cfnResource)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal cloudformation resource to yaml: %s", err))
	}
	textLines := common.GetLinesFromBytes(yamlBytes)

	return textLines
}

func findTagsLinesYAML(textLines []string) common.Lines {
	tagsLines := common.Lines{Start: -1, End: len(textLines) - 1}
	indent := ""
	for i, line := range textLines {
		if strings.Contains(line, TagsAttributeName+":") {
			tagsLines.Start = i
			indent = common.ExtractIndentationOfLine(line)
		} else if common.ExtractIndentationOfLine(line) <= indent && !strings.HasPrefix(line, indent+"-") && tagsLines.Start >= 0 {
			tagsLines.End = i - 1
			return tagsLines
		}
	}

	return tagsLines
}
