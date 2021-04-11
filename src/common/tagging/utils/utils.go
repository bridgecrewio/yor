package utils

import (
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/code2cloud"
	"bridgecrewio/yor/src/common/tagging/gittag"
	"bridgecrewio/yor/src/common/tagging/simple"
	"sort"
)

type TagGroupName string

const (
	SimpleTagGroupName TagGroupName = "simple"
	GitTagGroupName    TagGroupName = "git"
	Code2Cloud         TagGroupName = "code2cloud"
)

var tagGroupsByName = map[TagGroupName]tagging.ITagGroup{
	SimpleTagGroupName: &simple.TagGroup{},
	GitTagGroupName:    &gittag.TagGroup{},
	Code2Cloud:         &code2cloud.TagGroup{},
}

func TagGroupsByName(name TagGroupName) tagging.ITagGroup {
	tagGroup, ok := tagGroupsByName[name]
	if !ok {
		return nil
	}

	return tagGroup
}

func GetAllTagGroupsNames() []string {
	tagGroupNames := make([]string, 0)
	for name := range tagGroupsByName {
		tagGroupNames = append(tagGroupNames, string(name))
	}
	sort.Strings(tagGroupNames)

	return tagGroupNames
}

func ListAllTagsKeys() []string {
	allTagsKeys := make([]string, 0)
	for _, tagGroup := range tagGroupsByName {
		defaultTags := tagGroup.GetDefaultTags()
		for _, tag := range defaultTags {
			tag.Init()
			allTagsKeys = append(allTagsKeys, tag.GetKey())
		}
	}

	return allTagsKeys
}
