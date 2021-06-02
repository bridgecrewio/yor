package gittag

import (
	"fmt"
	"reflect"

	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type GitLastModifiedAtTag struct {
	tags.Tag
}

func (t *GitLastModifiedAtTag) Init() {
	t.Key = gitLastModifiedAtTagKey
}

func (t *GitLastModifiedAtTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	latestCommit := gitBlame.GetLatestCommit()
	if latestCommit == nil {
		return nil, fmt.Errorf("latest commit is unavailable")
	}
	return &tags.Tag{Key: t.Key, Value: latestCommit.Date.UTC().Format("2006-01-02 15:04:05")}, nil
}

func (t *GitLastModifiedAtTag) GetDescription() string {
	return "The last time this resource's configuration was modified"
}
