package utils

import (
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/types"
	"bufio"
	"fmt"
	"os"
	"reflect"
)

type YamlParser struct {
	types.YamlParser
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

func (p *YamlParser) InSlice(slice interface{}, elem interface{}) bool {
	return InSlice(slice, elem)
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

func (p *YamlParser) StructContainsProperty(s interface{}, property string) (bool, reflect.Value) {
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

func GetFileScanner(filePath string, nonFoundLines *structure.Lines) (*bufio.Scanner, *structure.Lines) {
	//#nosec G304
	file, err := os.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
		return nil, nonFoundLines
	}
	scanner := bufio.NewScanner(file)
	defer func() {
		_ = file.Close()
	}()
	return scanner, nonFoundLines
}
