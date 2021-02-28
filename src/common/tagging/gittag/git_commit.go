package gittag

import (
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"reflect"
)

type GitCommitTag struct {
	tags.Tag
}

func (t *GitCommitTag) Init() {
	t.Key = "git_commit"
}

func (t *GitCommitTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}

	latestCommit := gitBlame.GetLatestCommit()
	if latestCommit == nil {
		return nil, fmt.Errorf("latest commit is unavailable")
	}
	return &tags.Tag{Key: t.Key, Value: latestCommit.Hash.String()}, nil
}

func (t *GitCommitTag) GetDescription() string {
	return "Commit tag (i.e. d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0)"
}
