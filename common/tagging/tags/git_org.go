package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"fmt"
	"reflect"
)

type GitOrgTag struct {
	Tag
}

func (t *GitOrgTag) Init() ITag {
	t.Key = "git_org"

	return t
}

func (t *GitOrgTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	t.Value = gitBlame.GitOrg
	return nil
}
