package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"fmt"
	"reflect"
)

type GitLastModifiedAtTag struct {
	Tag
}

func (t *GitLastModifiedAtTag) Init() {
	t.Key = "git_last_modified_at"
}

func (t *GitLastModifiedAtTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	t.Value = gitBlame.GetLatestCommit().Date.String()
	return nil
}
