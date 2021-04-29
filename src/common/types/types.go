package types

import (
	"bridgecrewio/yor/src/common"

	"github.com/awslabs/goformation/v4/cloudformation"
)

type Block struct {
	common.Block
}

type YamlParser struct {
	RootDir              string
	FileToResourcesLines map[string]common.Lines
}

type ServerlessTemplate struct {
	Service  string `yaml:"service"`
	Provider struct {
		Name         string            `yaml:"name"`
		Runtime      string            `yaml:"runtime"`
		Region       string            `yaml:"region"`
		ProviderTags map[string]string `yaml:"tags"`
		CFNTags      map[string]string `yaml:"stackTags"`
		Functions    interface{}       `yaml:"functions"`
		Resources    struct {
			Resources *cloudformation.Template `yaml:"Resources"`
		} `yaml:"resources"`
	} `yaml:"provider"`
}

type ServerlessParser struct {
	YamlParser YamlParser
	Template   *ServerlessTemplate
}
