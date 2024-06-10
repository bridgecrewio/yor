package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
)

// RemoveGcpInvalidChars Source of regex: https://cloud.google.com/compute/docs/labeling-resources
var RemoveGcpInvalidChars = regexp.MustCompile(`[^\p{Ll}\p{Lo}\p{N}_-]`)
var SkipResourcesByComment = make([]string, 0)
var mutex sync.Mutex

func AppendSkipedByCommentToRunnerSkippedResources(skippedResources *[]string) {
	mutex.Lock()
	*skippedResources = append(*skippedResources, SkipResourcesByComment...)
	SkipResourcesByComment = SkipResourcesByComment[:0]
	mutex.Unlock()
}

func InSlice[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func SliceInSlices[T comparable](elems [][]T, vSlice []T) bool {
	for _, elemSlice := range elems {
		curSize := len(elemSlice)
		if curSize != len(vSlice) {
			continue
		}
		equalNum := 0
		for i, elem := range elemSlice {
			if elem == vSlice[i] {
				equalNum++
			}
		}
		if equalNum == curSize {
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
		content, _ := os.ReadFile(absFilePath)
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

func MaxMapCountKey(m map[string]int) string {
	var maxKey string
	var maxCount = -1
	for key, count := range m {
		if count > maxCount {
			maxKey = key
			maxCount = count
		}
	}
	return maxKey
}
