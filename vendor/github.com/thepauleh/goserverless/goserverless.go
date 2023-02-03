package goserverless

import (
	"encoding/json"
	"io/ioutil"

	"github.com/awslabs/goformation/v4/intrinsics"

	"github.com/thepauleh/goserverless/serverless"
)

// Open and parse a Serverless template from file.
// Works with YAML formatted templates.
func Open(filename string) (*serverless.Template, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return openYaml(data)
}

func openYaml(input []byte) (*serverless.Template, error) {
	intrinsified, err := intrinsics.ProcessYAML(input, nil)
	if err != nil {
		return nil, err
	}
	template := &serverless.Template{}
	if err := json.Unmarshal(intrinsified, template); err != nil {
		return nil, err
	}

	return template, nil
}
