package utils

import (
	"bridgecrewio/yor/src/common"
	"reflect"
)

type YamlParser struct {
	RootDir              string
	FileToResourcesLines map[string]common.Lines
}

func (p *YamlParser) InSlice(slice interface{}, elem interface{}) bool {
	for _, e := range p.convertToInterfaceSlice(slice) {
		if p.getKind(e) != p.getKind(elem) {
			continue
		}
		if p.getKind(e) == reflect.Slice {
			inSlice := true
			for _, subElem := range p.convertToInterfaceSlice(elem) {
				inSlice = inSlice && p.InSlice(e, subElem)
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

func (p *YamlParser) getKind(val interface{}) reflect.Kind {
	s := reflect.ValueOf(val)
	return s.Kind()
}

func (p *YamlParser) convertToInterfaceSlice(origin interface{}) []interface{} {
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
