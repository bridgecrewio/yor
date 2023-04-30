package slsParser

import "github.com/bridgecrewio/goformation/v5/cloudformation"

// Template definition - only important fields
type Template struct {
	FrameworkVersion string                  `json:"frameworkVersion,omitempty"`
	Functions        map[string]Function     `json:"functions,omitempty"`
	Resources        cloudformation.Template `json:"resources,omitempty"`
}

// Function definition - only important fields
type Function struct {
	Name string                 `json:"name,omitempty"`
	Tags map[string]interface{} `json:"tags,omitempty"`
}
