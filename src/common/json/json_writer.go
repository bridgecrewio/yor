package json

import (
	"encoding/json"
	"fmt"
	"github.com/bridgecrewio/yor/src/common/types"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/utils"
)

// readFilePath: origin file path
// blocks: resources blocks, some containing updated tags
// writeFilePath: destination for writing the updated json
// fileBracketsPairs: mapping of all brackets in file by their start char index
// The function updates the content of `readFilePath` with updated tags from `blocks` and writes it to `writeFilePath`
func WriteJsonFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, fileBracketsPairs map[int]types.BracketPair) error {

	// #nosec G304
	originFileSrc, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}
	originFileStr := string(originFileSrc)

	newStringsByStartChar := make(map[int]string) // map between start char index and the string that should be written in that index
	Start2EndCharMap := make(map[int]int)         // map start index to end index
	for _, resourceBlock := range blocks {
		if resourceBlock.IsBlockTaggable() {
			tagsDiff := resourceBlock.CalculateTagsDiff()
			if len(tagsDiff.Added) == 0 && len(tagsDiff.Updated) == 0 {
				// if resource was not changed during the run, continue
				continue
			}

			resourceBrackets := FindScopeInJSON(originFileStr, resourceBlock.GetResourceID(), fileBracketsPairs, &structure.Lines{Start: -1, End: -1})
			Start2EndCharMap[resourceBrackets.Open.CharIndex] = resourceBrackets.Close.CharIndex
			newResourceLines := AddTagsToResourceStr(originFileStr, resourceBlock, fileBracketsPairs)
			newStringsByStartChar[resourceBrackets.Open.CharIndex] = newResourceLines
		}
	}

	// write changes
	textToWrite := originFileStr
	if len(newStringsByStartChar) > 0 {
		// sort start chars in ascending order
		startChars := make([]int, 0, len(newStringsByStartChar))
		for c := range newStringsByStartChar {
			startChars = append(startChars, c)
		}
		sort.Sort(sort.IntSlice(startChars))

		textToWrite = ""
		lastReplacedIndex := 0
		for _, cIndex := range startChars {
			// write text until next changed string and append new string
			textToWrite += originFileStr[lastReplacedIndex:cIndex] + newStringsByStartChar[cIndex]
			// set the pointer of the string continuation to be the end index of the replaced part.
			lastReplacedIndex = Start2EndCharMap[cIndex] + 1
		}
		textToWrite += originFileStr[lastReplacedIndex:]
	}

	err = ioutil.WriteFile(writeFilePath, []byte(textToWrite), 0600)
	return err
}

// The function gets the entire context as a string, and returns a string of a resource with the updated tags
func AddTagsToResourceStr(fullOriginStr string, resourceBlock structure.IBlock, fileBracketsPairs map[int]types.BracketPair) string {
	logger.Debug(fmt.Sprintf("setting tags to resource %s in path %s", resourceBlock.GetResourceID(), resourceBlock.GetFilePath()))
	updatedTags := resourceBlock.MergeTags()

	// extract the resource's brackets scope and get the origin str for that resource
	resourceBrackets := FindScopeInJSON(fullOriginStr, resourceBlock.GetResourceID(), fileBracketsPairs, &structure.Lines{Start: -1, End: -1})
	resourceStr := fullOriginStr[resourceBrackets.Open.CharIndex : resourceBrackets.Close.CharIndex+1]

	tagsAttributeName := resourceBlock.GetTagsAttributeName()
	indexOfTags := findJSONKeyIndex(resourceStr, tagsAttributeName) // get the start index of the tags key in the resource string

	if indexOfTags >= 0 {
		// extract the tags' brackets scope and get the origin str for them
		tagBrackets := FindScopeInJSON(fullOriginStr, tagsAttributeName, fileBracketsPairs, &structure.Lines{Start: resourceBrackets.Open.Line, End: resourceBrackets.Close.Line})
		tagsStr := fullOriginStr[tagBrackets.Open.CharIndex:tagBrackets.Close.CharIndex]

		//	now find the indentation of the first tags entry by searching an indent between "[" and first "{". If there is a newline, restart the indent.
		tagBlockIndent := findIndent(tagsStr, '{', 0)                               // find the indent of each tag block " { "
		tagEntryIndent := findIndent(tagsStr, '"', strings.Index(tagsStr[1:], "{")) // find the indent of the key and value entry

		// unmarshal updated tags with the indent matching origin file
		strUpdatedTags, err := json.MarshalIndent(updatedTags, tagBlockIndent, strings.TrimPrefix(tagEntryIndent, tagBlockIndent))
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to unmarshal tags %s with indent '%s' because of error: %s", updatedTags, tagBlockIndent, err))
		}
		tagsStartRelativeToResource := tagBrackets.Open.CharIndex - resourceBrackets.Open.CharIndex
		tagsEndRelativeToResource := tagBrackets.Close.CharIndex - resourceBrackets.Open.CharIndex

		// set the resource string with the updated and indented tags
		resourceStr = resourceStr[:tagsStartRelativeToResource] + string(strUpdatedTags) + resourceStr[tagsEndRelativeToResource+1:]
	} else {
		// step 1 - extract the parent of the tags attribute from the new resource (not from the file)
		jsonResourceStr := getJSONStr(resourceBlock.GetRawBlock()) // encode raw block to json
		identifiersToAdd := make([]string, 0)
		parentIdentifier := tagsAttributeName

		// step 2 - find the parent identifier in the origin resource. If not found continue to look for identifiers until reaching the resource name
		indexOfParent := -1
		for indexOfParent < 0 && parentIdentifier != resourceBlock.GetResourceID() {
			identifiersToAdd = append(identifiersToAdd, parentIdentifier)
			parentIdentifier = FindParentIdentifier(jsonResourceStr, parentIdentifier)
			indexOfParent = findJSONKeyIndex(resourceStr, parentIdentifier)
		}

		// step 3 - find indent from last parent scope start to it's first child
		topIdentifierScope := FindScopeInJSON(fullOriginStr, identifiersToAdd[len(identifiersToAdd)-1], fileBracketsPairs, &structure.Lines{Start: resourceBrackets.Open.Line, End: resourceBrackets.Close.Line})
		indent := findIndent(fullOriginStr, '"', topIdentifierScope.Open.CharIndex)

		// step 4 - add the missing data

		// create a map of data to add
		entriesToAdd := make(map[string]interface{})
		for i := len(identifiersToAdd) - 1; i <= 0; i++ {
			if i > 0 {
				entriesToAdd[identifiersToAdd[i]] = make(map[string]interface{})
			} else {
				entriesToAdd[identifiersToAdd[i]] = updatedTags
			}
		}

		// marshal the map using the extracted indentation
		jsonToAdd, err := json.MarshalIndent(entriesToAdd, indent, "\t")
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to unmarshal tags %s with indent '%s' because of error: %s", entriesToAdd, indent, err))
		}
		textToAdd := string(jsonToAdd)

		// remove first and last chars, which are '{' and '}' - we already have the top level map and don't need it
		textToAdd = textToAdd[1 : len(textToAdd)-1]

		// fix indentation after removing the top level map
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

		// add comma after tags
		textToAdd = "\n" + strings.Join(editedLines, "\n") + ","

		// adding the tags as the first element
		resourceStr = resourceStr[:(topIdentifierScope.Open.CharIndex+1)-resourceBrackets.Open.CharIndex] + textToAdd + resourceStr[(topIdentifierScope.Open.CharIndex+1)-resourceBrackets.Open.CharIndex:]
	}

	return resourceStr
}

// find the indentation in a string `str` from starting char index until `charToStop` is identified
// if a newline is encountered, restart the indentation to ""
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

// marshal an interface into json and return a string of that json
func getJSONStr(rawBlock interface{}) string {
	jsonBytes, err := json.Marshal(rawBlock)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal resource to json: %s", err))
	}

	return string(jsonBytes)
}

// map the lines of all resources in a file and return it with the brackets mapping
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
		matchingBrackets := FindScopeInJSON(fileStr, resourceName, bracketPairs, &structure.Lines{Start: -1, End: -1})
		resourceToLines[resourceName] = &structure.Lines{Start: matchingBrackets.Open.Line, End: matchingBrackets.Close.Line}
	}

	return resourceToLines, bracketPairs
}

// find all brackets in a string
func MapBracketsInString(str string) []types.Brackets {
	allBrackets := make([]types.Brackets, 0)
	lineCounter := 1
	for cIndex, c := range str {
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

// given a list of all brackets of pair, map all the pairs correctly and return them ordered by the open char index
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
func FindScopeInJSON(str string, key string, bracketsPairs map[int]types.BracketPair, linesRange *structure.Lines) types.BracketPair {
	indexOfKey := -1
	if linesRange.Start != -1 {
		fileLines := strings.Split(str, "\n")
		beforeRange := strings.Join(fileLines[:linesRange.Start], "\n")
		rangeLinesStr := strings.Join(fileLines[linesRange.Start:linesRange.End], "\n")
		indexOfKey = findJSONKeyIndex(rangeLinesStr, key)
		indexOfKey = len(beforeRange) + indexOfKey
	} else {
		indexOfKey = findJSONKeyIndex(str, key)
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

// given a brackets pair, find the pair that wraps them
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

// find the identifier of the parent of a given child.
// for example, str = {parent: {child: [...] }} and childIdentifier="child", return "parent"
func FindParentIdentifier(str string, childIdentifier string) string {
	// create mapping of all brackets in resource
	bracketsInResourceJSON := MapBracketsInString(str)
	bracketsPairsInResourceJSON := GetBracketsPairs(bracketsInResourceJSON)

	// get tags brackets
	childScope := FindScopeInJSON(str, childIdentifier, bracketsPairsInResourceJSON, &structure.Lines{Start: -1, End: -1})
	// find the brackets that wrap the "tags"
	wrappingBracketsScope := FindWrappingBrackets(bracketsPairsInResourceJSON, childScope)
	// extract the name of the tags' parent (for example, in CFN it will be "Properties")
	r, _ := regexp.Compile("\"")
	quoteMarksIndexes := r.FindAllStringIndex(str[:wrappingBracketsScope.Open.CharIndex], -1)
	indexOfLastQuoteMark := quoteMarksIndexes[len(quoteMarksIndexes)-1][0]
	indexOfSecondToLastQuoteMark := quoteMarksIndexes[len(quoteMarksIndexes)-2][0]
	parentIdentifier := str[indexOfSecondToLastQuoteMark+1 : indexOfLastQuoteMark]

	return parentIdentifier
}

// find the index of an entry in a JSON by wrapping it with "<key>":
func findJSONKeyIndex(str string, key string) int {
	r, _ := regexp.Compile("\"" + key + `"\s*:`) // support a case of one or more spaces before colon
	indexPair := r.FindStringIndex(str)
	if len(indexPair) == 0 {
		return -1
	}

	return indexPair[0]
}
