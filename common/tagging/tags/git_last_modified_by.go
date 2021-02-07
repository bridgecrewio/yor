package tags

import (
	"bridgecrewio/yor/common/git_service"
	"fmt"
	"reflect"
	"time"
)

type GitLastModifiedByTag struct {
	Tag
}

func (t *GitLastModifiedByTag) Init() ITag {
	t.Key = "git_last_modified_by"

	return t
}

func (t *GitLastModifiedByTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*git_service.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	latestDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	modifyingUser := ""
	for _, v := range gitBlame.BlamesByLine {
		if latestDate.Before(v.Date) {
			latestDate = v.Date
			modifyingUser = v.Author
		}
	}

	t.Value = modifyingUser
	return nil
}
