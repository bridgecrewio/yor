package gittag

import (
	"fmt"
	"reflect"

	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type GitFileTag struct {
	tags.Tag
}

func (t *GitFileTag) Init() {
	t.Key = tags.GitFileTagKey
}

func (t *GitFileTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	return &tags.Tag{Key: t.Key, Value: gitBlame.FilePath}, nil
}

func (t *GitFileTag) GetDescription() string {
	return "The file (including path) in the repository where this resource is provisioned in IaC"
}
