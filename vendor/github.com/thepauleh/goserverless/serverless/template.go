package serverless

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/sanathkr/yaml"
)

// Template represents an AWS CloudFormation template
// see: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-anatomy.html
type Template struct {
	Service          *Service                `json:"service,omitempty"`
	FrameworkVersion string                  `json:"frameworkVersion,omitempty"`
	Functions        map[string]Function     `json:"functions,omitempty"`
	Resources        cloudformation.Template `json:"resources,omitempty"`
	Provider         *Provider               `json:"provider,omitempty"`
	Package          *Package                `json:"package,omitempty"`
	Plugins          []string                `json:"plugins,omitempty"`
	Custom           map[string]interface{}  `json:"custom,omitempty"`
}

// NewTemplate creates a new AWS CloudFormation template struct
func NewTemplate(serviceName string) *Template {
	return &Template{
		Service:          &Service{Name: serviceName},
		FrameworkVersion: ">=1.0.0 <2.0.0",
		Functions:        map[string]Function{},
	}
}

// YAML converts an AWS CloudFormation template object to YAML
func (t *Template) YAML() ([]byte, error) {
	return yaml.Marshal(t)
}
