package gittag

import (
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"reflect"
)

type GitRepoTag struct {
	tags.Tag
}

func (t *GitRepoTag) Init() {
	t.Key = "git_repo"
}

func (t *GitRepoTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	return &tags.Tag{Key: t.Key, Value: gitBlame.GitRepository}, nil
}

func (t *GitRepoTag) GetDescription() string {
	return "The repository where this resource is provisioned in IaC"
}
