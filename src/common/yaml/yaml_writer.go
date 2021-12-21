package yaml

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/utils"
	"github.com/sanathkr/yaml"
	"github.com/thepauleh/goserverless/serverless"
)

const SingleIndent = "  "

func WriteYAMLFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, tagsAttributeName string, resourcesStartToken string) error {
	// #nosec G304
	// read file bytes
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}
	isCfn := !strings.Contains(filepath.Base(readFilePath), "serverless")
	originLines := utils.GetLinesFromBytes(originFileSrc)

	oldResourcesLineRange := computeResourcesLineRange(originLines, blocks, isCfn)
	resourcesLines := make([]string, 0)
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].GetLines().Start < blocks[j].GetLines().Start
	})
	linesPerTag := 1
	if isCfn {
		linesPerTag = 2
	}
	for _, resourceBlock := range blocks {
		rawBlock := resourceBlock.GetRawBlock()
		newResourceLines := getYAMLLines(rawBlock, isCfn)
		newResourceTagLineRange, _ := FindTagsLinesYAML(newResourceLines, tagsAttributeName)
		oldResourceLinesRange := resourceBlock.GetLines()
		oldResourceLines := originLines[oldResourceLinesRange.Start : oldResourceLinesRange.End+1]

		// if the block is not taggable, write it and continue
		if !resourceBlock.IsBlockTaggable() {
			resourcesLines = append(resourcesLines, oldResourceLines...)
			continue
		}

		oldResourceTagLines := resourceBlock.GetTagsLines()
		// if the resource doesn't contain Tags entry - create it
		if oldResourceTagLines.Start == -1 || oldResourceTagLines.End == -1 {
			// get the indentation of the property under the resource name
			tagAttributeIndent := ExtractIndentationOfLine(oldResourceLines[1])
			if isCfn {
				tagAttributeIndent += SingleIndent
			}
			lastIndex := -1
			for i, line := range oldResourceLines {
				if len(ExtractIndentationOfLine(line)) < len(tagAttributeIndent) {
					continue
				}
				lastIndex = i
			}
			resourcesLines = append(resourcesLines, oldResourceLines[:lastIndex+1]...)
			resourcesLines = append(resourcesLines, tagAttributeIndent+tagsAttributeName+":") // add the 'Tags:' line
			tagIndent := tagAttributeIndent
			if isCfn {
				tagIndent += SingleIndent
			}
			resourcesLines = append(resourcesLines, IndentLines(newResourceLines[newResourceTagLineRange.Start+1:newResourceTagLineRange.End+1], tagIndent)...)
			resourcesLines = append(resourcesLines, oldResourceLines[lastIndex+1:]...)
			continue
		}

		oldTagsIndent := ExtractIndentationOfLine(oldResourceLines[oldResourceTagLines.Start-oldResourceLinesRange.Start])
		if isCfn {
			oldTagsIndent += SingleIndent
		}
		resourcesLines = append(resourcesLines, oldResourceLines[:oldResourceTagLines.Start-oldResourceLinesRange.Start]...) // add all the resource's line before the tags
		tagLines := oldResourceLines[oldResourceTagLines.Start-oldResourceLinesRange.Start : oldResourceTagLines.End-oldResourceLinesRange.Start+1]
		diff := resourceBlock.CalculateTagsDiff()
		if isCfn {
			UpdateExistingCFNTags(tagLines, diff.Updated)
		} else {
			UpdateExistingSLSTags(tagLines, diff.Updated)
		}
		allNewResourceTagLines := IndentLines(newResourceLines[newResourceTagLineRange.Start+1:newResourceTagLineRange.End+1], oldTagsIndent)
		var netNewResourceLines []string
		for i := 0; i < len(allNewResourceTagLines); i += linesPerTag {
			l := allNewResourceTagLines[i]
			key := getKeyFromLine(l, isCfn)
			if key == "" {
				continue
			}
			found := false
			for _, tag := range diff.Added {
				if tag.GetKey() == key {
					found = true
					break
				}
			}
			if found {
				netNewResourceLines = append(netNewResourceLines, allNewResourceTagLines[i:i+linesPerTag]...)
			}
		}
		resourcesLines = append(resourcesLines, tagLines...)            // Add old tags
		resourcesLines = append(resourcesLines, netNewResourceLines...) // Add new tags
		// Add any other attributes after the tags
		resourcesLines = append(resourcesLines, oldResourceLines[oldResourceTagLines.End-oldResourceLinesRange.Start+1:]...)
	}
	allLines := make([]string, oldResourcesLineRange.Start)
	copy(allLines, originLines[:oldResourcesLineRange.Start])
	if !isCfn {
		allLines = append(allLines, resourcesStartToken+":")
	}
	allLines = append(allLines, resourcesLines...)
	allLines = append(allLines, originLines[oldResourcesLineRange.End+1:]...)
	linesText := strings.Join(allLines, "\n")

	err = ioutil.WriteFile(writeFilePath, []byte(linesText), 0600)

	return err
}

func getKeyFromLine(l string, isCfn bool) string {
	if isCfn {
		if strings.Contains(l, " Key:") {
			return strings.ReplaceAll(strings.ReplaceAll(l, " ", ""), "-Key:", "")
		}
	} else {
		return strings.Split(strings.ReplaceAll(strings.ReplaceAll(l, " ", ""), "-", ""), ":")[0]
	}
	return ""
}

func UpdateExistingCFNTags(tagsLinesList []string, diff []*tags.TagDiff) {
	currentValueLine := -1
	valueToSet := ""

	for i, tagLine := range tagsLinesList {
		if strings.Contains(tagLine, ` Key:`) {
			for _, tag := range diff {
				keyr := regexp.MustCompile(`\b` + tag.Key + `\b`)
				if keyr.Match([]byte(tagLine)) {
					if currentValueLine > -1 {
						tagsLinesList[currentValueLine] = ReplaceTagValue(tagsLinesList[currentValueLine], tag.NewValue)
						currentValueLine = -1
					} else {
						valueToSet = tag.NewValue
					}
					continue
				}
			}
		}
		if strings.Contains(tagLine, ` Value:`) {
			if valueToSet != "" {
				tagsLinesList[i] = ReplaceTagValue(tagLine, valueToSet)
				valueToSet = ""
			} else {
				currentValueLine = i
			}
		}
	}
}

func ReplaceTagValue(line string, value string) string {
	tr := regexp.MustCompile(`\bValue\s*:\s*.*`)
	return tr.ReplaceAllString(line, `Value: `+value)
}

func UpdateExistingSLSTags(tagLines []string, diff []*tags.TagDiff) {
	for i, line := range tagLines {
		key := strings.Split(strings.ReplaceAll(line, " ", ""), ":")[0]
		for _, tag := range diff {
			if key == tag.Key {
				lineWithoutValue := strings.Split(line, ":")[0]
				tagLines[i] = lineWithoutValue + ": " + tag.NewValue
			}
		}
	}
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

func getYAMLLines(rawBlock interface{}, isCfn bool) []string {
	var textLines []string
	yamlBytes, err := yaml.Marshal(rawBlock)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal resource to yaml: %s", err))
	}

	textLines = utils.GetLinesFromBytes(yamlBytes)

	if !isCfn {
		slsFunction := rawBlock.(serverless.Function)
		if utils.AllNil(slsFunction.VPC.SecurityGroupIds, slsFunction.VPC.SubnetIds) {
			textLines = removeLineByAttribute(textLines, "vpc:")
		}
		if utils.AllNil(slsFunction.Package.Include, slsFunction.Package.Artifact, slsFunction.Package.Exclude,
			slsFunction.Package.ExcludeDevDependencies, slsFunction.Package.Individually) {
			textLines = removeLineByAttribute(textLines, "package")
		}
	}
	return textLines
}

func removeLineByAttribute(textLines []string, attribute string) []string {
	vpcLineIndex := -1
	for i, line := range textLines {
		if strings.Contains(line, attribute) {
			vpcLineIndex = i
			break
		}
	}
	if vpcLineIndex != -1 {
		textLines = append(textLines[:vpcLineIndex], textLines[vpcLineIndex+1:]...)
	}
	return textLines
}

func FindTagsLinesYAML(textLines []string, tagsAttributeName string) (structure.Lines, bool) {
	tagsLines := structure.Lines{Start: -1, End: -1}
	var lineIndent string
	var tagsExist bool
	var tagsIndent = ""
	for i, line := range textLines {
		lineIndent = ExtractIndentationOfLine(line)
		switch {
		case strings.HasPrefix(strings.TrimSpace(line), tagsAttributeName+":"):
			tagsLines.Start = i
			tagsIndent = lineIndent
			tagsExist = true
		case lineIndent <= tagsIndent && (tagsLines.Start >= 0 || i == len(textLines)-1):
			tagsLines.End = findLastNonEmptyLine(textLines, i-1)
			return tagsLines, tagsExist
		case i == len(textLines)-1 && !tagsExist:
			return tagsLines, tagsExist
		}
	}
	if !tagsExist {
		tagsLines.Start = tagsLines.End
	} else if tagsLines.End == -1 {
		tagsLines.End = findLastNonEmptyLine(textLines, len(textLines)-1)
	}
	return tagsLines, tagsExist
}

func MapResourcesLineYAML(filePath string, resourceNames []string, resourcesStartToken string) map[string]*structure.Lines {
	resourceToLines := make(map[string]*structure.Lines)
	for _, resourceName := range resourceNames {
		// initialize a map between resource name and its lines in file
		resourceToLines[resourceName] = &structure.Lines{Start: -1, End: -1}
	}
	// #nosec G304
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
		return nil
	}

	readResources := false
	latestResourceName := ""
	fileLines := strings.Split(string(file), "\n")
	resourcesIndent := 0
	// iterate file line by line
	for i, line := range fileLines {
		cleanContent := strings.TrimSpace(line)
		if strings.HasPrefix(cleanContent, resourcesStartToken+":") {
			readResources = true
			resourcesIndent = countLeadingSpaces(line)
			continue
		}

		if readResources {
			lineIndent := countLeadingSpaces(line)
			if lineIndent <= resourcesIndent && strings.TrimSpace(line) != "" && !strings.Contains(line, "#") {
				// No longer inside resources block, get the last line of the previous resource if exists
				readResources = false
				if latestResourceName != "" {
					resourceToLines[latestResourceName].End = findLastNonEmptyLine(fileLines, i-1)
				}
				break
			}
			for _, resName := range resourceNames {
				resNameRegex := regexp.MustCompile(fmt.Sprintf("^ {1,5}%v:", resName))
				if resNameRegex.Match([]byte(line)) {
					if latestResourceName != "" {
						// Complete previous function block
						resourceToLines[latestResourceName].End = findLastNonEmptyLine(fileLines, i-1)
					}
					latestResourceName = resName
					resourceToLines[latestResourceName].Start = i
					break
				}
			}
			if !strings.HasPrefix(line, " ") && strings.TrimSpace(line) != "" && readResources && latestResourceName != "" && !strings.HasPrefix(line, "#") {
				// This is no longer in the functions block, complete last function block
				resourceToLines[latestResourceName].End = findLastNonEmptyLine(fileLines, i-1)
				break
			}
		}
	}
	if resourceToLines[latestResourceName].End == -1 {
		// Handle last line of resource is last line of file
		resourceToLines[latestResourceName].End = findLastNonEmptyLine(fileLines, len(fileLines)-1)
	}
	return resourceToLines
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func findLastNonEmptyLine(fileLines []string, maxIndex int) int {
	for i := utils.MinInt(maxIndex, len(fileLines)-1); i >= 0; i-- {
		if strings.TrimSpace(fileLines[i]) != "" {
			return i
		}
	}
	return 0
}

func IndentLines(textLines []string, indent string) []string {
	for i, originLine := range textLines {
		noLeadingWhitespace := strings.TrimLeft(originLine, "\t \n")
		if strings.Contains(originLine, "- Key") {
			textLines[i] = indent + noLeadingWhitespace
		} else {
			textLines[i] = indent + SingleIndent + noLeadingWhitespace
		}
	}

	return textLines
}

func ExtractIndentationOfLine(textLine string) string {
	indent := ""
	for _, c := range textLine {
		if c != ' ' && c != '-' {
			break
		}
		indent += " "
	}

	return indent
}
