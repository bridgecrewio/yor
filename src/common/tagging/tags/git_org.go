package tags

import (
	"bridgecrewio/yor/src/common/gitservice"
	"fmt"
	"reflect"
)

type GitOrgTag struct {
	Tag
}

func (t *GitOrgTag) Init() {
	t.Key = "git_org"
}

func (t *GitOrgTag) CalculateValue(data interface{}) (ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	return &Tag{Key: t.Key, Value: gitBlame.GitOrg}, nil
}
