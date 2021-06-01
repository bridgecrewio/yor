package utils

import (
	"sort"

	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/code2cloud"
	"github.com/bridgecrewio/yor/src/common/tagging/external"
	"github.com/bridgecrewio/yor/src/common/tagging/gittag"
	"github.com/bridgecrewio/yor/src/common/tagging/simple"
)

type TagGroupName string

const (
	SimpleTagGroupName TagGroupName = "simple"
	GitTagGroupName    TagGroupName = "git"
	Code2Cloud         TagGroupName = "code2cloud"
	ExternalTagName    TagGroupName = "external"
)

var tagGroupsByName = map[TagGroupName]tagging.ITagGroup{
	SimpleTagGroupName: &simple.TagGroup{},
	GitTagGroupName:    &gittag.TagGroup{},
	Code2Cloud:         &code2cloud.TagGroup{},
}

func TagGroupsByName(name TagGroupName) tagging.ITagGroup {
	var tagGroup tagging.ITagGroup
	switch name {
	case SimpleTagGroupName:
		tagGroup = &simple.TagGroup{}
	case GitTagGroupName:
		tagGroup = &gittag.TagGroup{}
	case Code2Cloud:
		tagGroup = &code2cloud.TagGroup{}
	case ExternalTagName:
		tagGroup = &external.TagGroup{}
	}

	return tagGroup
}

func GetAllTagGroupsNames() []string {
	tagGroupNames := make([]string, 0)
	for name := range tagGroupsByName {
		tagGroupNames = append(tagGroupNames, string(name))
	}
	sort.Strings(tagGroupNames)
	tagGroupNames = append(tagGroupNames, string(ExternalTagName)) // Add the external tag name as the last tag group
	return tagGroupNames
}
