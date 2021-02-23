package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"errors"
	"fmt"
	"reflect"
)

type GitCommitTag struct {
	Tag
}

func (t *GitCommitTag) Init() {
	t.Key = "git_commit"
}

func (t *GitCommitTag) CalculateValue(data interface{}) (ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	latestCommit := gitBlame.GetLatestCommit()
	if latestCommit == nil {
		return nil, errors.New("latest commit is unavailable")
	}
	return &Tag{Key: t.Key, Value: latestCommit.Hash.String()}, nil
}
