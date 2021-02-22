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

func (t *GitLastModifiedByTag) CalculateValue(data interface{}) (ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	return &Tag{Key: t.Key, Value: gitBlame.GetLatestCommit().Author}, nil
}
