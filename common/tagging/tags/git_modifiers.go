package tags

import (
	"bridgecrewio/yor/common/git_service"
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

func (t *GitModifiersTag) CalculateValue(data interface{}) error {
	gitBlame, ok := data.(*git_service.GitBlame)
	if !ok {
		return fmt.Errorf("failed to convert data to *GitBlame, which is required to calculte tag value. Type of data: %s", reflect.TypeOf(data))
	}
	var modifyingUsers []string
	for _, v := range gitBlame.BlamesByLine {
		userName := strings.Split(v.Author, "@")[0]
		modifyingUsers = append(modifyingUsers, userName)
	}

	sort.Strings(modifyingUsers)

	t.Value = strings.Join(modifyingUsers, "/")
	return nil
}
