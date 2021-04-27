package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/sanathkr/yaml"
)

func WriteYAMLFile(readFilePath string, blocks []IBlock, writeFilePath string, resourcesLinesRange common.Lines, tagsAttributeName string) error {
	// read file bytes
	// #nosec G304
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}
	originLines := common.GetLinesFromBytes(originFileSrc)

	resourcesIndent := common.ExtractIndentationOfLine(originLines[resourcesLinesRange.Start])

	resourcesLines := make([]string, 0)
	for _, resourceBlock := range blocks {
		newResourceLines := getYAMLLines(resourceBlock.GetRawBlock(), tagsAttributeName)
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
			tagAttributeIndent := common.ExtractIndentationOfLine(oldResourceLines[1])
			resourcesLines = append(resourcesLines, oldResourceLines...)                      // add all the existing resource data first
			resourcesLines = append(resourcesLines, tagAttributeIndent+tagsAttributeName+":") // add the 'Tags:' line
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

	err = ioutil.WriteFile(writeFilePath, []byte(linesText), 0600)

	return err
}
func reflectValueToMap(rawMap reflect.Value, currResourceMap *map[string]interface{}, tagsAttributeName string) *map[string]interface{} {
	for _, mapKeyRef := range rawMap.MapKeys() {
		mapKey := mapKeyRef.Elem().String()
		mapValue := rawMap.MapIndex(mapKeyRef)
		if mapKey == tagsAttributeName {
			fmt.Println(1)
		} else {
			switch mapValue.Elem().Kind() {
			case reflect.String:
				(*currResourceMap)[mapKey] = mapValue.Elem().String()
				break
			case reflect.Int:
				(*currResourceMap)[mapKey] = mapValue.Elem().Int()
				break
			}
		}
	}
	return currResourceMap
}

func getYAMLLines(rawBlock interface{}, tagsAttributeName string) []string {
	var textLines []string
	var yamlBytes []byte
	var err error
	switch rawBlock.(type) {
	case reflect.Value:
		newResourceMapTemp := make(map[string]interface{}, 0)
		newResourceMap := reflectValueToMap(rawBlock.(reflect.Value), &newResourceMapTemp, tagsAttributeName)
		yamlBytes, err = yaml.Marshal(newResourceMap)
		break
	default:
		yamlBytes, err = yaml.Marshal(rawBlock)

	}
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal resource to yaml: %s", err))
	}
	textLines = common.GetLinesFromBytes(yamlBytes)

	return textLines
}

func FindTagsLinesYAML(textLines []string, tagsAttributeName string) common.Lines {
	tagsLines := common.Lines{Start: -1, End: len(textLines) - 1}
	indent := ""
	for i, line := range textLines {
		if strings.Contains(line, tagsAttributeName+":") {
			tagsLines.Start = i + 1
			indent = common.ExtractIndentationOfLine(line)
		} else if common.ExtractIndentationOfLine(line) <= indent && (tagsLines.Start >= 0 || i == len(textLines)-1) {
			tagsLines.End = i
			return tagsLines
		}
	}

	return tagsLines
}
