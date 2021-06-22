package gittag

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type GitModifiersTag struct {
	tags.Tag
}

func (t *GitModifiersTag) Init() {
	t.Key = tags.GitModifiersTagKey
}

func (t *GitModifiersTag) CalculateValue(data interface{}) (tags.ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	foundModifyingUsers := make(map[string]bool)
	var modifyingUsers []string
	for _, v := range gitBlame.BlamesByLine {
		userName := strings.Split(v.Author, "@")[0]
		if !foundModifyingUsers[userName] && userName != "" && !strings.Contains(userName, "[") {
			modifyingUsers = append(modifyingUsers, userName)
			foundModifyingUsers[userName] = true
		}
	}

	sort.Strings(modifyingUsers)

	return &tags.Tag{Key: t.Key, Value: strings.Join(modifyingUsers, "/")}, nil
}

func (t *GitModifiersTag) GetDescription() string {
	return "The users who modified this resource"
}
