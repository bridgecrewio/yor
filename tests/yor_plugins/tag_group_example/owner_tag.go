package main

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

type DirTag struct {
	tags.Tag
}

func (d *DirTag) Init() {
	d.Key = "custom_owner"
}

func (d *DirTag) CalculateValue(block interface{}) (tags.ITag, error) {
	blockVal, ok := block.(common.IBlock)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to IBlock, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(block))
	}

	dir := filepath.Dir(blockVal.GetFilePath())
	var owner string
	switch {
	case strings.HasPrefix(dir, "src/auth"):
		owner = "team-infra@company.com"
	case strings.HasPrefix(dir, "data/"):
		owner = "team-data@company.com"
	case strings.HasPrefix(dir, "jenkins"):
		owner = "team-ops@company.com"
	default:
		owner = "team-it@company.com"
	}

	return &tags.Tag{Key: d.Key, Value: owner}, nil
}
