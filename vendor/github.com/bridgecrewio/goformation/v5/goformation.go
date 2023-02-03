package goformation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sanathkr/yaml"

	"github.com/bridgecrewio/goformation/v5/cloudformation"
	"github.com/bridgecrewio/goformation/v5/intrinsics"
)

//go:generate generate/generate.sh

// Open and parse a AWS CloudFormation template from file.
// Works with either JSON or YAML formatted templates.
func Open(filename string) (*cloudformation.Template, error) {
	return OpenWithOptions(filename, nil)
}

// OpenWithOptions opens and parse a AWS CloudFormation template from file.
// Works with either JSON or YAML formatted templates.
// Parsing can be tweaked via the specified options.
func OpenWithOptions(filename string, options *intrinsics.ProcessorOptions) (*cloudformation.Template, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filename, ".json") {
		// This is definitely JSON
		return ParseJSONWithOptions(data, options)
	}

	return ParseYAMLWithOptions(data, options)
}

// StringifyInnerYAMLValues converts values in YAML to string, for a keyPath stated.
// keyPath should look like []string{"Resources", "*", "Properties", "Environment", "Variables", "*"}
func StringifyInnerYAMLValues(iMapData interface{}, keyPath []string) interface{} {
	if len(keyPath) == 0 {
		return fmt.Sprintf("%v", iMapData)
	}

	mapData, isMap := iMapData.(map[string]interface{})

	currProperty := keyPath[0]

	if currProperty == "*" {
		if isMap {
			for key, val := range mapData {
				mapData[key] = StringifyInnerYAMLValues(val, keyPath[1:])
			}
			return mapData
		}

		sliceData, isSlice := iMapData.([]interface{})
		if !isSlice {
			return iMapData
		}
		for i, val := range sliceData {
			sliceData[i] = StringifyInnerYAMLValues(val, keyPath[1:])
		}
		return sliceData
	}

	// if the current property is not '*', i.e a named property
	if !isMap {
		return iMapData
	}
	innerMap, propertyFound := mapData[currProperty]
	if propertyFound {
		mapData[currProperty] = StringifyInnerYAMLValues(innerMap, keyPath[1:])
	}

	return mapData
}

// StringifyYAMLValues converts values in YAML to string, for all keyPaths stated.
// keyPaths should look like []string{"Resources/*/Properties/Environment/Variables/*"}
func StringifyYAMLValues(data []byte, keyPaths []string) ([]byte, error) {
	var structData interface{}
	err := yaml.Unmarshal(data, &structData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	mapData, ok := structData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert sturct to map for stringifying YAML elements")
	}

	for _, keyPath := range keyPaths {
		iMapData := StringifyInnerYAMLValues(mapData, strings.Split(keyPath, "/"))
		mapData, ok = iMapData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to stringify paths in YAML")
		}
	}

	updatedData, err := yaml.Marshal(&mapData)

	return updatedData, err
}

// ParseYAMLWithOptions an AWS CloudFormation template (expects a []byte of valid YAML)
// Parsing can be tweaked via the specified options.
func ParseYAMLWithOptions(data []byte, options *intrinsics.ProcessorOptions) (*cloudformation.Template, error) {
	var err error
	if options != nil && len(options.StringifyPaths) > 0 {
		data, err = StringifyYAMLValues(data, options.StringifyPaths)
		if err != nil {
			return nil, err
		}
	}

	// Process all AWS CloudFormation intrinsic functions (e.g. Fn::Join)
	intrinsified, err := intrinsics.ProcessYAML(data, options)
	if err != nil {
		return nil, err
	}

	return unmarshal(intrinsified)

}

// ParseJSON an AWS CloudFormation template (expects a []byte of valid JSON)
func ParseJSON(data []byte) (*cloudformation.Template, error) {
	return ParseJSONWithOptions(data, nil)
}

// ParseJSONWithOptions an AWS CloudFormation template (expects a []byte of valid JSON)
// Parsing can be tweaked via the specified options.
func ParseJSONWithOptions(data []byte, options *intrinsics.ProcessorOptions) (*cloudformation.Template, error) {

	// Process all AWS CloudFormation intrinsic functions (e.g. Fn::Join)
	intrinsified, err := intrinsics.ProcessJSON(data, options)
	if err != nil {
		return nil, err
	}

	return unmarshal(intrinsified)

}

func unmarshal(data []byte) (*cloudformation.Template, error) {

	template := &cloudformation.Template{}
	if err := json.Unmarshal(data, template); err != nil {
		return nil, err
	}

	return template, nil

}
