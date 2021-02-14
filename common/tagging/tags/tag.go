package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"time"

	"github.com/go-git/go-git/v5"
)

type Tag struct {
	Key   string
	Value string
}

var TagTypes = []ITag{
	&GitOrgTag{},
	&GitRepoTag{},
	&GitFileTag{},
	&GitCommitTag{},
	&GitModifiersTag{},
	&GitLastModifiedAtTag{},
	&GitLastModifiedByTag{},
}

type ITag interface {
	Init() ITag
	CalculateValue(data interface{}) error
	GetKey() string
	GetValue() string
}

type TagDiff struct {
	Key       string
	PrevValue string
	NewValue  string
}

func Init(key string, value string) ITag {
	return &Tag{
		Key:   key,
		Value: value,
	}
}

func (t *Tag) Init() ITag {
	return t
}

func (t *Tag) CalculateValue(_ interface{}) error {
	return nil
}

func (t *Tag) GetKey() string {
	return t.Key
}

func (t *Tag) GetValue() string {
	return t.Value
}

func getLatestCommit(blame *gitservice.GitBlame) (latestCommit *git.Line) {
	latestDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, v := range blame.BlamesByLine {
		if latestDate.Before(v.Date) {
			latestDate = v.Date
			latestCommit = v
		}
	}
	return
}
