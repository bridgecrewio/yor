package slsParser

import (
	"encoding/json"
	"github.com/bridgecrewio/goformation/v5/intrinsics"
	"os"
)

// Open and parse a Serverless template from file.
// Works with YAML formatted templates.
func Open(filename string) (*Template, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return openYaml(data)
}

func openYaml(input []byte) (*Template, error) {
	intrinsified, err := intrinsics.ProcessYAML(input, nil)
	if err != nil {
		return nil, err
	}
	template := &Template{}
	if err := json.Unmarshal(intrinsified, template); err != nil {
		return nil, err
	}

	return template, nil
}
