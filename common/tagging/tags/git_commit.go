package tags

import (
	"bridgecrewio/yor/common/git_service"
	"fmt"
	"reflect"
)

type GitCommitTag struct {
	Tag
}

func (t *GitCommitTag) Init() ITag {
	t.Key = "git_commit"

	return t
}

func (t *GitCommitTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*git_service.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	t.Value = GetLatestCommit(gitBlame).Hash.String()
	return nil
}
