package main

import (
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"path/filepath"
	"reflect"
)

type DirTag struct {
	tags.Tag
}

func (d *DirTag) Init() {
	d.Key = "bc_dir"
}

func (d *DirTag) CalculateValue(block interface{}) (tags.ITag, error) {
	blockVal, ok := block.(structure.IBlock)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to IBlock, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(block))
	}

	dir := filepath.Dir(blockVal.GetFilePath())

	return &tags.Tag{Key: d.Key, Value: dir}, nil
}
