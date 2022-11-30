package gittag

import (
	"fmt"
	"reflect"

	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type GitLastModifiedByTag struct {
	tags.Tag
}

func (t *GitLastModifiedByTag) Init() {
	t.Key = tags.GitLastModifiedByTagKey
}

func (t *GitLastModifiedByTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	latestCommit := gitBlame.GetLatestCommit()
	if latestCommit == nil {
		return nil, fmt.Errorf("latest commit is unavailable")
	}
	return &tags.Tag{Key: t.Key, Value: latestCommit.Author}, nil
}

func (t *GitLastModifiedByTag) GetDescription() string {
	return "The last user who modified this resource"
}
