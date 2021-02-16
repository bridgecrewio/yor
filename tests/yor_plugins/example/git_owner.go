package main

import (
	"bridgecrewio/yor/common/gitservice"
	"fmt"
	"reflect"
)

type GitOwnerTag struct {
	Key   string
	Value string
}

func (t *GitOwnerTag) Init() {
	t.Key = "git_owner"
}

func (t *GitOwnerTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	t.Value = gitBlame.GetLatestCommit().Author

	return nil
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
