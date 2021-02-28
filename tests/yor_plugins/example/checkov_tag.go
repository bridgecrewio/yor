package main

import "bridgecrewio/yor/src/common/tagging/tags"

type CheckovTag struct {
	Key   string
	Value string
}

func (t *CheckovTag) GetDescription() string {
	return "Checkov tag - adds the static tag {yor_checkov:checkov} to every taggable resource"
}

func (t *CheckovTag) Init() {
	t.Key = "yor_checkov"
}

func (t *CheckovTag) CalculateValue(_ interface{}) (tags.ITag, error) {
	return &tags.Tag{Key: t.Key, Value: "checkov"}, nil
}

func (t *CheckovTag) GetKey() string {
	return t.Key
}

func (t *CheckovTag) GetValue() string {
	return t.Value
}

func (t *CheckovTag) GetPriority() int {
	return 1
}
