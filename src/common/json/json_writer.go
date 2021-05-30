package json

import (
	"encoding/json"
	"fmt"
	"github.com/bridgecrewio/yor/src/common/types"
	"io/ioutil"
	"math"
	"sort"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/utils"
)

func WriteJsonFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, fileBracketsPairs map[int]types.BracketPair) error {
	// #nosec G304
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}
	originFileStr := string(originFileSrc)

	newStringsByStartChar := make(map[int]string)
	Start2EndCharMap := make(map[int]int)
	for _, resourceBlock := range blocks {
		if resourceBlock.IsBlockTaggable() {
			tagsDiff := resourceBlock.CalculateTagsDiff()
			if len(tagsDiff.Added) == 0 && len(tagsDiff.Updated) == 0 {
				// if resource was not changed during the run, continue
				continue
			}

			resourceBrackets := FindScopeInJSON(originFileStr, JsonEntry(resourceBlock.GetResourceID()), fileBracketsPairs, &structure.Lines{Start: -1, End: -1})
			Start2EndCharMap[resourceBrackets.Open.CharIndex] = resourceBrackets.Close.CharIndex
			newResourceLines := AddTagsToResourceStr(originFileStr, resourceBlock, fileBracketsPairs)
			newStringsByStartChar[resourceBrackets.Open.CharIndex] = newResourceLines
		}
	}

	textToWrite := originFileStr

	if len(newStringsByStartChar) > 0 {
		startChars := make([]int, 0, len(newStringsByStartChar))
		for c := range newStringsByStartChar {
			startChars = append(startChars, c)
		}

		// sort end chars in descending order
		sort.Sort(sort.IntSlice(startChars))
		// take all the text after the last resource to edit
		textToWrite = ""

		lastReplaced := 0
		for _, cIndex := range startChars {
			textToWrite = textToWrite + originFileStr[lastReplaced:cIndex] + newStringsByStartChar[cIndex]
			lastReplaced = Start2EndCharMap[cIndex] + 1
		}
		textToWrite = textToWrite + originFileStr[lastReplaced:]
	}

	err = ioutil.WriteFile(writeFilePath, []byte(textToWrite), 0600)
	return err
}

func AddTagsToResourceStr(originFileStr string, resourceBlock structure.IBlock, fileBracketsPairs map[int]types.BracketPair) string {
	updatedTags := resourceBlock.MergeTags()
	resourceBrackets := FindScopeInJSON(originFileStr, JsonEntry(resourceBlock.GetResourceID()), fileBracketsPairs, &structure.Lines{Start: -1, End: -1})
	resourceStr := originFileStr[resourceBrackets.Open.CharIndex : resourceBrackets.Close.CharIndex+1]

	tagsAttributeName := resourceBlock.GetTagsAttributeName()
	indexOfTags := strings.Index(resourceStr, JsonEntry(tagsAttributeName))

	if indexOfTags >= 0 {
		//	tags exists in resource
		tagBrackets := FindScopeInJSON(originFileStr, JsonEntry(tagsAttributeName), fileBracketsPairs, &structure.Lines{Start: resourceBrackets.Open.Line, End: resourceBrackets.Close.Line})
		//	now find the indentation of the first tags entry by searching an indent between "[" and first "{". If there is a newline, restart the indent.
		tagsStr := originFileStr[tagBrackets.Open.CharIndex:tagBrackets.Close.CharIndex]
		tagBlockIndent := findIndent(tagsStr, '{', 0)                                                                              // find the indent of each tag block " { "
		tagEntryIndent := findIndent(tagsStr, '"', strings.Index(tagsStr[1:], "{"))                                                // find the indent of the key and value entry
		strUpdatedTags, err := json.MarshalIndent(updatedTags, tagBlockIndent, strings.TrimPrefix(tagEntryIndent, tagBlockIndent)) // unmarshal updated tags with the indent matching origin file
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to unmarshal tags %s with indent '%s' because of error: %s", updatedTags, tagBlockIndent, err))
		}
		tagsStartRelativeToResource := tagBrackets.Open.CharIndex - resourceBrackets.Open.CharIndex
		tagsEndRelativeToResource := tagBrackets.Close.CharIndex - resourceBrackets.Open.CharIndex

		resourceStr = resourceStr[:tagsStartRelativeToResource] + string(strUpdatedTags) + resourceStr[tagsEndRelativeToResource+1:]
	} else {
		//	tags doesnt exist and we need to add it
		// step 1 - extract the parent of the tags attriubte
		jsonResourceStr := getJSONStr(resourceBlock.GetRawBlock()) // encode raw block to json
		identifiersToAdd := []string{}
		parentIdentifier := tagsAttributeName

		// step 2 - find the parent identifier in the origin resource. If not found continue to look for identifiers until reaching the resource name
		indexOfParent := -1
		for indexOfParent < 0 && parentIdentifier != resourceBlock.GetResourceID() {
			identifiersToAdd = append(identifiersToAdd, parentIdentifier)
			parentIdentifier = FindParentIdentifier(jsonResourceStr, parentIdentifier)
			indexOfParent = strings.Index(resourceStr, parentIdentifier)
		}

		// step 3 - find indent from last parent scope start to it's first child
		topIdentifierScope := FindScopeInJSON(originFileStr, identifiersToAdd[len(identifiersToAdd)-1], fileBracketsPairs, &structure.Lines{Start: resourceBrackets.Open.Line, End: resourceBrackets.Close.Line})
		indent := findIndent(originFileStr, '"', topIdentifierScope.Open.CharIndex)

		// step 4 - add the missing data
		entriesToAdd := make(map[string]interface{})
		//var entriesToAdd interface{}
		for i := len(identifiersToAdd) - 1; i <= 0; i++ {
			if i > 0 {
				entriesToAdd[identifiersToAdd[i]] = make(map[string]interface{})
			} else {
				entriesToAdd[identifiersToAdd[i]] = updatedTags
			}
		}

		jsonToAdd, err := json.MarshalIndent(entriesToAdd, indent, "\t")
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to unmarshal tags %s with indent '%s' because of error: %s", entriesToAdd, indent, err))
		}

		textToAdd := string(jsonToAdd)
		textToAdd = textToAdd[1 : len(textToAdd)-1]
		lines := strings.Split(textToAdd, "\n")
		editedLines := make([]string, 0)
		for _, l := range lines {
			for c := range l {
				if !utils.IsCharWhitespace(l[c]) {
					newL := strings.Replace(l, "\t", "", 1)
					editedLines = append(editedLines, newL)
					break
				}
			}
		}
		textToAdd = "\n" + strings.Join(editedLines, "\n") + "," // remove open and close brackets
		resourceStr = resourceStr[:(topIdentifierScope.Open.CharIndex+1)-resourceBrackets.Open.CharIndex] + textToAdd + resourceStr[(topIdentifierScope.Open.CharIndex+1)-resourceBrackets.Open.CharIndex:]
	}

	return resourceStr
}

func findIndent(str string, charToStop byte, startIndex int) string {
	indent := ""
	charIndex := startIndex
	currChar := str[charIndex]
	for currChar != charToStop {
		if utils.IsCharWhitespace(currChar) {
			if currChar == '\n' {
				indent = ""
			} else {
				indent += string(currChar)
			}
		} else {
			indent = ""
		}
		charIndex++
		currChar = str[charIndex]
	}

	return indent
}

func getJSONStr(rawBlock interface{}) string {
	jsonBytes, err := json.Marshal(rawBlock)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal resource to json: %s", err))
	}

	return string(jsonBytes)
}

func MapResourcesLineJSON(filePath string, resourceNames []string) (map[string]*structure.Lines, map[int]types.BracketPair) {
	resourceToLines := make(map[string]*structure.Lines)
	// #nosec G304
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
		return nil, nil
	}

	fileStr := string(file)
	bracketsInFile := MapBracketsInString(fileStr)
	bracketPairs := GetBracketsPairs(bracketsInFile)

	for _, resourceName := range resourceNames {
		matchingBrackets := FindScopeInJSON(fileStr, JsonEntry(resourceName), bracketPairs, &structure.Lines{Start: -1, End: -1})
		resourceToLines[resourceName] = &structure.Lines{Start: matchingBrackets.Open.Line, End: matchingBrackets.Close.Line}
	}

	return resourceToLines, bracketPairs
}

func JsonEntry(s string) string {
	return "\"" + s + "\""
}

func MapBracketsInString(fileStr string) []types.Brackets {
	allBrackets := make([]types.Brackets, 0)
	lineCounter := 1
	for cIndex, c := range fileStr {
		switch c {
		case '{':
			allBrackets = append(allBrackets, types.Brackets{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: lineCounter, CharIndex: cIndex})
		case '}':
			allBrackets = append(allBrackets, types.Brackets{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: lineCounter, CharIndex: cIndex})
		case '[':
			allBrackets = append(allBrackets, types.Brackets{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: lineCounter, CharIndex: cIndex})
		case ']':
			allBrackets = append(allBrackets, types.Brackets{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: lineCounter, CharIndex: cIndex})
		case '\n':
			lineCounter += 1
		}
	}

	return allBrackets
}

func GetBracketsPairs(bracketsInString []types.Brackets) map[int]types.BracketPair {
	startCharToBrackets := make(map[int]types.BracketPair)
	bracketShape2BracketsStacks := make(map[types.BracketShape][]types.Brackets)

	for _, bracket := range bracketsInString {
		stack, ok := bracketShape2BracketsStacks[bracket.Shape]
		if bracket.Type == types.OpenBrackets {
			if !ok {
				stack = make([]types.Brackets, 0)
			}
			stack = append(stack, bracket)
			bracketShape2BracketsStacks[bracket.Shape] = stack
		} else {
			if !ok {
				logger.Error("malformed json file", "SILENT")
			}
			openBracket := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			bracketShape2BracketsStacks[bracket.Shape] = stack
			startCharToBrackets[openBracket.CharIndex] = types.BracketPair{Open: openBracket, Close: bracket}
		}
	}

	return startCharToBrackets
}

// Find the index of a key in json string and return the start and end brackets of the key scope
func FindScopeInJSON(fileStr string, key string, bracketsPairs map[int]types.BracketPair, linesRange *structure.Lines) types.BracketPair {
	var indexOfKey int
	if linesRange.Start != -1 {
		fileLines := strings.Split(fileStr, "\n")
		beforeRange := strings.Join(fileLines[:linesRange.Start], "\n")
		rangeLinesStr := strings.Join(fileLines[linesRange.Start:linesRange.End], "\n")
		indexOfKey = strings.Index(rangeLinesStr, key)
		indexOfKey = len(beforeRange) + indexOfKey
	} else {
		indexOfKey = strings.Index(fileStr, key)
	}

	foundBracketPair := types.BracketPair{}

	nextIndex := math.MaxInt64
	for index, bracketPair := range bracketsPairs {
		if index > indexOfKey {
			if bracketPair.Open.CharIndex < nextIndex {
				nextIndex = bracketPair.Open.CharIndex
				foundBracketPair = bracketPair
			}
		}
	}

	return foundBracketPair
}

func FindWrappingBrackets(allBracketPairs map[int]types.BracketPair, innerBracketPair types.BracketPair) types.BracketPair {
	wrappingPair := -1
	for i, bracketsPair := range allBracketPairs {
		if bracketsPair.Open.CharIndex < innerBracketPair.Open.CharIndex && bracketsPair.Close.CharIndex > innerBracketPair.Close.CharIndex {
			if wrappingPair == -1 || (bracketsPair.Open.CharIndex > allBracketPairs[wrappingPair].Open.CharIndex && bracketsPair.Close.CharIndex < allBracketPairs[wrappingPair].Close.CharIndex) {
				wrappingPair = i
			}
		}
	}

	return allBracketPairs[wrappingPair]
}

func FindParentIdentifier(str string, childIdentifier string) string {
	// create mapping of all brackets in resource
	bracketsInResourceJSON := MapBracketsInString(str)
	bracketsPairsInResourceJSON := GetBracketsPairs(bracketsInResourceJSON)

	// get tags brackets
	tagsScope := FindScopeInJSON(str, childIdentifier, bracketsPairsInResourceJSON, &structure.Lines{Start: -1, End: -1})
	// find the brackets that wrap the "tags"
	wrappingBracketsScope := FindWrappingBrackets(bracketsPairsInResourceJSON, tagsScope)
	// extract the name of the tags' parent (for example, in CFN it will be "Properties")
	indexOfLastQuoteMark := strings.LastIndex(str[:wrappingBracketsScope.Open.CharIndex], "\"")
	indexOfSecondToLastQuoteMark := strings.LastIndex(str[:indexOfLastQuoteMark], "\"")
	parentIdentifier := str[indexOfSecondToLastQuoteMark+1 : indexOfLastQuoteMark]

	return parentIdentifier
}
