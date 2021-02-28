package gittag

import (
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"reflect"
)

type GitLastModifiedAtTag struct {
	tags.Tag
}

func (t *GitLastModifiedAtTag) Init() {
	t.Key = "git_last_modified_at"
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
	return "Last Modified At - the last time this resource configuration was updated"
}
