package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"fmt"
	"reflect"
)

type GitLastModifiedByTag struct {
	Tag
}

func (t *GitLastModifiedByTag) Init() {
	t.Key = "git_last_modified_by"
}

func (t *GitLastModifiedByTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	t.Value = getLatestCommit(gitBlame).Author
	return nil
}
