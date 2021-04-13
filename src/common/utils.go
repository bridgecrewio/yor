package common

import (
	"io/ioutil"
	"reflect"
	"strings"
)

type Lines struct {
	Start int
	End   int
}

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

func GetFileFormat(filePath string) string {
	splitByDot := strings.Split(filePath, ".")
	if len(splitByDot) < 2 {
		return ""
	}
	if strings.HasSuffix(filePath, CFTFileType.Extension) {
		content, _ := ioutil.ReadFile(filePath)
		if strings.HasPrefix(string(content), "{") {
			return JSONFileType.FileFormat
		}
		return YamlFileType.FileFormat
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
