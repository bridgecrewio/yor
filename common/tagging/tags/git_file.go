package tags

import (
	"bridgecrewio/yor/common/git_service"
	"fmt"
	"reflect"
)

type GitFileTag struct {
	Tag
}

func (t *GitFileTag) Init() {
	t.Key = "git_file"
}

func (t *GitFileTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*git_service.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	t.Value = gitBlame.FilePath
	return nil
}
