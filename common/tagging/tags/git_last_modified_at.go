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

func (t *GitLastModifiedAtTag) CalculateValue(data interface{}) (ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	return &Tag{Key: t.Key, Value: gitBlame.GetLatestCommit().Date.UTC().Format("2006-01-02 15:04:05")}, nil
}
