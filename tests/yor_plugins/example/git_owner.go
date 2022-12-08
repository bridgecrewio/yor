package main

import (
	"fmt"
	"reflect"

	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type GitOwnerTag struct {
	Key   string
	Value string
}

func (t *GitOwnerTag) GetDescription() string {
	return "Tag that marks the organizational owner of the resource"
}

func (t *GitOwnerTag) SetTagPrefix(s string) {
	t.Key = s + t.Key
}

func (t *GitOwnerTag) Init() {
	t.Key = "git_owner"
}

func (t *GitOwnerTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	return &tags.Tag{Key: t.Key, Value: gitBlame.GetLatestCommit().Author}, nil
}

func (t *GitOwnerTag) GetKey() string {
	return t.Key
}

func (t *GitOwnerTag) GetValue() string {
	return t.Value
}

func (t *GitOwnerTag) GetPriority() int {
	return -1
}

func (t *GitOwnerTag) SetValue(val string) {
	t.Value = val
}
