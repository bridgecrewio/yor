package yaml

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/utils"

	"github.com/sanathkr/yaml"
)

func WriteYAMLFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, resourcesLinesRange structure.Lines, tagsAttributeName string) error {
	// read file bytes
	// #nosec G304
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}
	isCfn := !strings.Contains(filepath.Base(readFilePath), "serverless")
	originLines := utils.GetLinesFromBytes(originFileSrc)
	originLines = utils.ReorderByTags(originLines, tagsAttributeName, isCfn)
	resourcesIndent := utils.ExtractIndentationOfLine(originLines[resourcesLinesRange.Start])

	oldResourcesLineRange := computeResourcesLineRange(originLines, blocks, isCfn)
	resourcesLines := make([]string, 0)
	for _, resourceBlock := range blocks {
		rawBlock := resourceBlock.GetRawBlock()
		var newResourceLines []string
		switch rawBlockCasted := rawBlock.(type) {
		case cloudformation.Resource:
			newResourceLines = getYAMLLines(rawBlockCasted, tagsAttributeName, isCfn)
		case map[interface{}]interface{}:
			newResourceLines = getYAMLLines(resourceBlock.GetRawBlock().(map[interface{}]interface{}), tagsAttributeName, isCfn)
		}
		newResourceTagLineRange, _ := FindTagsLinesYAML(newResourceLines, tagsAttributeName)
		oldResourceLinesRange := resourceBlock.GetLines()
		oldResourceLines := originLines[oldResourceLinesRange.Start-1 : oldResourceLinesRange.End]

		// if the block is not taggable, write it and continue
		if !resourceBlock.IsBlockTaggable() {
			resourcesLines = append(resourcesLines, oldResourceLines...)
			continue
		}

		oldResourceTagLines, oldTagsExist := FindTagsLinesYAML(oldResourceLines, tagsAttributeName)
		// if the resource don't contain Tags entry - create it
		if oldResourceTagLines.Start == -1 || oldResourceTagLines.End == -1 {
			// get the indentation of the property under the resource name
			tagAttributeIndent := utils.ExtractIndentationOfLine(oldResourceLines[1])
			resourcesLines = append(resourcesLines, oldResourceLines...)                      // add all the existing resource data first
			resourcesLines = append(resourcesLines, tagAttributeIndent+tagsAttributeName+":") // add the 'Tags:' line
			// add the tags with extra indentation below the 'Tags:' line
			resourcesLines = append(resourcesLines, utils.IndentLines(newResourceLines[newResourceTagLineRange.Start+1:newResourceTagLineRange.End+1], tagAttributeIndent+"  ")...)
			continue
		}

		oldTagsIndent := utils.ExtractIndentationOfLine(oldResourceLines[oldResourceTagLines.Start])
		resourcesLines = append(resourcesLines, oldResourceLines[:oldResourceTagLines.Start]...) // add all the resource's line before the tags
		if !oldTagsExist {
			oldTagsIndent = resourcesIndent
			resourcesLines = append(resourcesLines, resourcesIndent+fmt.Sprintf("%s:", tagsAttributeName))
		}
		resourcesLines = append(resourcesLines, utils.IndentLines(newResourceLines[newResourceTagLineRange.Start:newResourceTagLineRange.End], oldTagsIndent)...) // add tags
		resourcesLines = append(resourcesLines, utils.IndentLines(newResourceLines[newResourceTagLineRange.End:], oldTagsIndent)...)
	}
	allLines := make([]string, len(originLines))
	if isCfn {
		copy(allLines, originLines[:oldResourcesLineRange.Start-1])
	} else {
		copy(allLines, originLines[:oldResourcesLineRange.Start+1])
	}
	allLines = append(allLines, resourcesLines...)
	if !isCfn {
		allLines = append(allLines, originLines[resourcesLinesRange.End:]...)
	}
	filteredLines := make([]string, 0)
	for i, line := range allLines {
		if strings.Trim(line, "\n \t") == "" {
			continue
		} else {
			filteredLines = append(filteredLines, allLines[i])
		}
	}
	linesText := strings.Join(filteredLines, "\n")

	err = ioutil.WriteFile(writeFilePath, []byte(linesText), 0600)

	return err
}

func computeResourcesLineRange(originLines []string, blocks []structure.IBlock, isCfn bool) structure.Lines {
	ret := structure.Lines{
		Start: -1,
		End:   -1,
	}
	minLine := math.Inf(0)
	maxLine := -1
	for _, block := range blocks {
		minLine = math.Min(minLine, float64(block.GetLines().Start))
		maxLine = int(math.Max(float64(maxLine), float64(block.GetLines().End)))
	}
	if !isCfn {
		functionsBlockStartLine := -1
		for i, line := range originLines {
			if strings.Contains(line, "functions:") {
				functionsBlockStartLine = i
				break
			}
		}
		minLine = math.Min(minLine, float64(functionsBlockStartLine))
	}
	ret.Start = int(minLine)
	ret.End = maxLine
	return ret
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

func getYAMLLines(rawBlock interface{}, tagsAttributeName string, isCfn bool) []string {
	var textLines []string
	var castedRawBlock = rawBlock
	tempTagsMap := make(map[string]interface{})
	_, ok := rawBlock.(map[interface{}]interface{})
	if ok {
		castedRawBlock = reflectValueToMap(rawBlock, &tempTagsMap, tagsAttributeName)
	}
	yamlBytes, err := yaml.Marshal(castedRawBlock)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal resource to yaml: %s", err))
	}

	textLines = utils.GetLinesFromBytes(yamlBytes)
	textLines = utils.ReorderByTags(textLines, tagsAttributeName, isCfn)
	return textLines
}

func FindTagsLinesYAML(textLines []string, tagsAttributeName string) (structure.Lines, bool) {
	tagsLines := structure.Lines{Start: -1, End: len(textLines) - 1}
	var prevLine string
	var lineIndent string
	var tagsExist bool
	var tagsIndent = ""
	for i, line := range textLines {
		lineIndent = utils.ExtractIndentationOfLine(line)
		switch {
		case strings.Contains(line, tagsAttributeName+":"):
			tagsLines.Start = i + 1
			tagsIndent = lineIndent
			tagsExist = true
		case lineIndent < tagsIndent && (tagsLines.Start >= 0 || i == len(textLines)-1):
			tagsLines.End = i - 1
			return tagsLines, tagsExist
		case i == len(textLines)-1 && !tagsExist:
			tagsLines.End = i
			tagsLines.Start = tagsLines.End
			tagsIndent = utils.ExtractIndentationOfLine(prevLine) //nolint:ineffassign,staticcheck
			return tagsLines, tagsExist
		}
		prevLine = line
	}
	if !tagsExist {
		tagsLines.Start = tagsLines.End
	}
	return tagsLines, tagsExist
}

func EncodeBlocksToYaml(readFilePath string, blocks []structure.IBlock) []structure.IBlock {
	fileFormat := utils.GetFileFormat(readFilePath)
	switch fileFormat {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		for _, block := range blocks {
			yamlBlock := block.(IYamlBlock)
			yamlBlock.UpdateTags()
		}
		return blocks
	default:
		logger.Warning(fmt.Sprintf("unsupported file type %s", fileFormat))
		return nil
	}

}
