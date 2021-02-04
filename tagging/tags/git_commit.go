package tags

import (
	"bridgecrewio/yor/git_service"
	"bridgecrewio/yor/tagging"
	"fmt"
	"reflect"
)

type GitCommitTag struct {
	tagging.Tag
}

func (t *GitCommitTag) Init() tagging.ITag {
	t.Key = "git_commit"

	return t
}

func (t *GitCommitTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*git_service.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	//TODO - remove print and implement logic
	fmt.Print(gitBlame)
	return nil
}
