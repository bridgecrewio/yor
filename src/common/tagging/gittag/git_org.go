package gittag

import (
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"reflect"
)

type GitOrgTag struct {
	tags.Tag
}

func (t *GitOrgTag) Init() {
	t.Key = "git_org"
}

func (t *GitOrgTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	return &tags.Tag{Key: t.Key, Value: gitBlame.GitOrg}, nil
}

func (t *GitOrgTag) GetDescription() string {
	return "The entity which owns the repository where this resource is provisioned in IaC"
}
