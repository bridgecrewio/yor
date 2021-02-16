package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"fmt"
	"reflect"
)

type GitRepoTag struct {
	Tag
}

func (t *GitRepoTag) Init() {
	t.Key = "git_repo"
}

func (t *GitRepoTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	t.Value = gitBlame.GitRepository
	return nil
}
