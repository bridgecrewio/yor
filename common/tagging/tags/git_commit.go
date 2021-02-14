package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"fmt"
	"reflect"
)

type GitCommitTag struct {
	Tag
}

func (t *GitCommitTag) Init() {
	t.Key = "git_commit"
}

func (t *GitCommitTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	t.Value = getLatestCommit(gitBlame).Hash.String()
	return nil
}
