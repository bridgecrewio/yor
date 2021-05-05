package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
)

func InSlice(slice interface{}, elem interface{}) bool {
	for _, e := range convertToInterfaceSlice(slice) {
		if getKind(e) != getKind(elem) {
			continue
		}
		if getKind(e) == reflect.Slice {
			inSlice := true
			for _, subElem := range convertToInterfaceSlice(elem) {
				inSlice = inSlice && InSlice(e, subElem)
				if !inSlice {
					break
				}
			}
			if inSlice {
				return true
			}
		} else if e == elem {
			return true
		}
	}

	return false
}

func getKind(val interface{}) reflect.Kind {
	s := reflect.ValueOf(val)
	return s.Kind()
}

func convertToInterfaceSlice(origin interface{}) []interface{} {
	s := reflect.ValueOf(origin)
	if s.Kind() != reflect.Slice {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func GetFileScanner(filePath string, nonFoundLines *structure.Lines) (*os.File, *bufio.Scanner, *structure.Lines) {
	//#nosec G304
	file, err := os.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
		return file, nil, nonFoundLines
	}
	scanner := bufio.NewScanner(file)
	return file, scanner, nonFoundLines
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

func ReorderByTags(lines []string, tagsAttributeName string, isCfn bool) []string {
	tagsOriginalStartLineInd := -1
	var processedTags = false
	var tagsIndent string
	var functionsIndent = ""
	var processedFunctions = false
	var sortedLines = make([]string, len(lines))
	for i, line := range lines {
		lineIndent := ExtractIndentationOfLine(line)
		switch {
		case strings.Contains(line, "functions:"):
			processedFunctions = true
			functionsIndent = lineIndent + "  "
			sortedLines[i] = lines[i]
		case strings.Contains(line, tagsAttributeName+":") || (!isCfn && processedFunctions):
			tagsOriginalStartLineInd = i
			tagsIndent = lineIndent
			sortedLines[i] = lines[i]
			processedTags = true
		case !isCfn && lineIndent <= functionsIndent:
			sortedLines[i] = lines[i]
			processedTags = false
		case lineIndent <= tagsIndent && !CfnTagLine(line) && processedTags && !strings.Contains(line, "stackTags"):
			sortedLines[i] = ""
			sortedLines = insert(sortedLines, tagsOriginalStartLineInd-1, line)
		case i == len(lines)-1 || lineIndent < tagsIndent:
			sortedLines[i] = lines[i]
		default:
			sortedLines[i] = lines[i]
		}
	}
	return sortedLines
}

func CfnTagLine(line string) bool {
	return strings.Contains(line, "- Key") || strings.Contains(line, "Value")
}

func IndentLines(textLines []string, indent string) []string {
	for i, originLine := range textLines {
		noLeadingWhitespace := strings.TrimLeft(originLine, "\t \n")
		if strings.Contains(originLine, "- Key") {
			textLines[i] = strings.Replace(indent, " ", "", 2) + noLeadingWhitespace
		} else {
			textLines[i] = indent + "  " + noLeadingWhitespace
		}
	}

	return textLines
}

func insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func StructContainsProperty(s interface{}, property string) (bool, reflect.Value) {
	var field reflect.Value
	sValue := reflect.ValueOf(s)

	// Check if the passed interface is a pointer
	if sValue.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		field = sValue.FieldByName(property)
	} else {
		// 'dereference' with Elem() and get the field by name
		field = sValue.Elem().FieldByName(property)
	}

	if !field.IsValid() {
		return false, field
	}

	return true, field
}
