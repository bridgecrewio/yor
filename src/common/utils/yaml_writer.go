package utils

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/sanathkr/yaml"
)

func WriteYAMLFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, resourcesLinesRange structure.Lines, tagsAttributeName string) error {
	// read file bytes
	// #nosec G304
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}
	originLines := GetLinesFromBytes(originFileSrc)

	resourcesIndent := ExtractIndentationOfLine(originLines[resourcesLinesRange.Start])

	resourcesLines := make([]string, 0)
	for _, resourceBlock := range blocks {
		newResourceLines := getYAMLLines(resourceBlock.GetRawBlock().(map[interface{}]interface{}), tagsAttributeName)
		newResourceTagLineRange := FindTagsLinesYAML(newResourceLines, tagsAttributeName)

		oldResourceLinesRange := resourceBlock.GetLines()
		oldResourceLines := originLines[oldResourceLinesRange.Start-1 : oldResourceLinesRange.End]

		// if the block is not taggable, write it and continue
		if !resourceBlock.IsBlockTaggable() {
			resourcesLines = append(resourcesLines, oldResourceLines...)
			continue
		}

		oldResourceTagLines := FindTagsLinesYAML(oldResourceLines, tagsAttributeName)
		// if the resource don't contain Tags entry - create it
		if oldResourceTagLines.Start == -1 || oldResourceTagLines.End == -1 {
			// get the indentation of the property under the resource name
			tagAttributeIndent := ExtractIndentationOfLine(oldResourceLines[1])
			resourcesLines = append(resourcesLines, oldResourceLines...)                      // add all the existing resource data first
			resourcesLines = append(resourcesLines, tagAttributeIndent+tagsAttributeName+":") // add the 'Tags:' line
			// add the tags with extra indentation below the 'Tags:' line
			resourcesLines = append(resourcesLines, IndentLines(newResourceLines[newResourceTagLineRange.Start+1:newResourceTagLineRange.End+1], tagAttributeIndent+"  ")...)
			continue
		}

		oldTagsIndent := ExtractIndentationOfLine(oldResourceLines[oldResourceTagLines.Start+1])

		resourcesLines = append(resourcesLines, oldResourceLines[:oldResourceTagLines.Start]...)                                                              // add all the resource's line before the tags
		resourcesLines = append(resourcesLines, resourcesIndent+newResourceLines[newResourceTagLineRange.Start])                                              // add the 'Tags:' line
		resourcesLines = append(resourcesLines, IndentLines(newResourceLines[newResourceTagLineRange.Start+1:newResourceTagLineRange.End], oldTagsIndent)...) // add tags
		resourcesLines = append(resourcesLines, oldResourceLines[oldResourceTagLines.End+1:]...)                                                              // add rest of resource lines
	}

	allLines := append(originLines[:resourcesLinesRange.Start-1], resourcesLines...)
	allLines = append(allLines, originLines[resourcesLinesRange.End:]...)
	linesText := strings.Join(allLines, "\n")

	err = ioutil.WriteFile(writeFilePath, []byte(linesText), 0600)

	return err
}

func reflectValueToMap(rawMap interface{}, currResourceMap *map[string]interface{}, tagsAttributeName string) *map[string]interface{} {
	switch rawMap := rawMap.(type) {
	case map[interface{}]interface{}:
		rawMapCasted := rawMap
		for mapKey, mapValue := range rawMapCasted {
			mapKeyString := mapKey.(string)
			if mapKey == tagsAttributeName {
				currResourceMap = reflectValueToMap(mapValue.(map[string]string), currResourceMap, tagsAttributeName)
			} else {
				switch mapValue := mapValue.(type) {
				case string, int, bool:
					(*currResourceMap)[mapKeyString] = mapValue
				}
			}
		}
	case map[string]string:
		if _, ok := (*currResourceMap)[tagsAttributeName]; !ok {
			(*currResourceMap)[tagsAttributeName] = make(map[string]string)
		}
		rawMapCasted := rawMap
		for mapKey, mapValue := range rawMapCasted {
			(*currResourceMap)[tagsAttributeName].(map[string]string)[mapKey] = mapValue
		}
	}
	return currResourceMap
}

func getYAMLLines(rawBlock map[interface{}]interface{}, tagsAttributeName string) []string {
	var textLines []string
	tempTagsMap := make(map[string]interface{})
	castedRawBlock := reflectValueToMap(rawBlock, &tempTagsMap, tagsAttributeName)
	yamlBytes, err := yaml.Marshal(castedRawBlock)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal resource to yaml: %s", err))
	}
	textLines = GetLinesFromBytes(yamlBytes)

	return textLines
}

func FindTagsLinesYAML(textLines []string, tagsAttributeName string) structure.Lines {
	tagsLines := structure.Lines{Start: -1, End: len(textLines) - 1}
	indent := ""
	for i, line := range textLines {
		if strings.Contains(line, tagsAttributeName+":") {
			tagsLines.Start = i + 1
			indent = ExtractIndentationOfLine(line)
		} else if line <= indent && (tagsLines.Start >= 0 || i == len(textLines)-1) {
			tagsLines.End = i
			return tagsLines
		}
	}

	return tagsLines
}

func GetFileFormat(filePath string) string {
	splitByDot := strings.Split(filePath, ".")
	if len(splitByDot) < 2 {
		return ""
	}
	if strings.HasSuffix(filePath, common.CFTFileType.Extension) {
		absFilePath, _ := filepath.Abs(filePath)
		// #nosec G304 - file is from user
		content, _ := ioutil.ReadFile(absFilePath)
		if strings.HasPrefix(string(content), "{") {
			return common.JSONFileType.FileFormat
		}
		return common.YamlFileType.FileFormat
	}
	return splitByDot[len(splitByDot)-1]
}

func GetLinesFromBytes(bytes []byte) []string {
	return strings.Split(string(bytes), "\n")
}

func ExtractIndentationOfLine(textLine string) string {
	indent := ""
	for _, c := range textLine {
		if c != ' ' {
			break
		}
		indent += " "
	}

	return indent
}

func IndentLines(textLines []string, indent string) []string {
	originIndent := ExtractIndentationOfLine(textLines[0])
	for i, originLine := range textLines {
		noLeadingWhitespace := originLine[len(originIndent):]
		textLines[i] = indent + noLeadingWhitespace
	}

	return textLines
}

func EncodeBlocksToYaml(readFilePath string, blocks []structure.IBlock, writeFilePath string, tagsAttributeName string, resourceLines structure.Lines) []structure.IBlock {
	fileFormat := GetFileFormat(readFilePath)
	switch fileFormat {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		for _, block := range blocks {
			block.UpdateTags()
		}
		return blocks
	default:
		logger.Warning(fmt.Sprintf("unsupported file type %s", fileFormat))
		return nil
	}

}
