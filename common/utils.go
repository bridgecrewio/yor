package common

import (
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
	"reflect"
	"regexp"
)

func InSlice(slice interface{}, elem interface{}) bool {
	for _, e := range convertToInterfaceSlice(slice) {
		if getKind(e) != getKind(elem) {
			continue
		}
		if getKind(e) == reflect.Slice {
			for _, subElem := range convertToInterfaceSlice(elem) {
				if InSlice(e, subElem) {
					return true
				}
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
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Try to match the tag's key name with a potentially quoted string
func IsTagKeyMatch(tag tags.ITag, keyName string) bool {
	match, _ := regexp.Match(fmt.Sprintf(`^"?%s"?$`, regexp.QuoteMeta(keyName)), []byte(tag.GetKey()))
	return match
}

func StructContainsProperty(s interface{}, property string) (bool, reflect.Value) {
	sValue := reflect.ValueOf(s)

	// Check if the passed interface is a pointer
	if sValue.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		sValue = reflect.New(reflect.TypeOf(s))
	}

	// 'dereference' with Elem() and get the field by name
	field := sValue.Elem().FieldByName(property)
	if !field.IsValid() {
		return false, field
	}

	return true, field
}
