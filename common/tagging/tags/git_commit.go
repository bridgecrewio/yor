package tags

import (
	"bridgecrewio/yor/common/git_service"
	"fmt"
	"reflect"
	"time"
)

type GitCommitTag struct {
	Tag
}

func (t *GitCommitTag) Init() ITag {
	t.Key = "git_commit"

	return t
}

func (t *GitCommitTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*git_service.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	minTime := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	var hash string
	for _, v := range gitBlame.BlamesByLine {
		if minTime.Before(v.Date) {
			minTime = v.Date
			hash = v.Hash.String()
		}
	}

	t.Value = hash
	return nil
}
