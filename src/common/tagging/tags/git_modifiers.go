package tags

import (
	"bridgecrewio/yor/src/common/gitservice"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type GitModifiersTag struct {
	Tag
}

func (t *GitModifiersTag) Init() {
	t.Key = "git_modifiers"
}

func (t *GitModifiersTag) CalculateValue(data interface{}) (ITag, error) {
	gitBlame, ok := data.(*gitservice.GitBlame)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	foundModifyingUsers := make(map[string]bool)
	var modifyingUsers []string
	for _, v := range gitBlame.BlamesByLine {
		userName := strings.Split(v.Author, "@")[0]
		if !foundModifyingUsers[userName] && userName != "" {
			modifyingUsers = append(modifyingUsers, userName)
			foundModifyingUsers[userName] = true
		}
	}

	sort.Strings(modifyingUsers)

	return &Tag{Key: t.Key, Value: strings.Join(modifyingUsers, "/")}, nil
}
