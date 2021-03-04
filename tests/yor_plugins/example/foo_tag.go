package main

import "bridgecrewio/yor/src/common/tagging/tags"

type FooTag struct {
	Key   string
	Value string
}

func (t *FooTag) Init() {
	t.Key = "yor_foo"
}

func (t *FooTag) CalculateValue(_ interface{}) (tags.ITag, error) {
	return &tags.Tag{Key: t.Key, Value: "foo"}, nil
}

func (t *FooTag) GetKey() string {
	return t.Key
}

func (t *FooTag) GetValue() string {
	return t.Value
}

func (t *FooTag) GetPriority() int {
	return 1
}
