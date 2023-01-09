package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
)

// RemoveGcpInvalidChars Source of regex: https://cloud.google.com/compute/docs/labeling-resources
var RemoveGcpInvalidChars = regexp.MustCompile(`[^\p{Ll}\p{Lo}\p{N}_-]`)

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

func AllNil(vv ...interface{}) bool {
	for _, v := range vv {
		if reflect.ValueOf(v).Kind() == reflect.Ptr && !reflect.ValueOf(v).IsNil() {
			return false
		}
		if reflect.ValueOf(v).Kind() == reflect.String && v != "" {
			return false
		}
		if reflect.ValueOf(v).Kind() == reflect.Slice && !reflect.ValueOf(v).IsNil() {
			return false
		}
	}
	return true
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

func GetFileScanner(filePath string, nonFoundLines *structure.Lines) (*bufio.Scanner, *structure.Lines) {
	//#nosec G304
	file, err := os.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
		return nil, nonFoundLines
	}
	scanner := bufio.NewScanner(file)
	return scanner, nonFoundLines
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

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IsCharWhitespace(c byte) bool {
	newStr := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, string(c))

	return newStr != string(c)
}

func SplitStringByComma(input []string) []string {
	var ans []string
	for _, i := range input {
		if strings.Contains(i, ",") {
			ans = append(ans, strings.Split(i, ",")...)
		} else {
			ans = append(ans, i)
		}
	}
	return ans
}

func FindSubMatchByGroup(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	if match == nil {
		return nil
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
