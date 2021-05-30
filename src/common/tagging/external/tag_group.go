package external

import (
	"errors"
	"io/ioutil"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"gopkg.in/yaml.v2"
)

type TagGroup struct {
	tagging.TagGroup
	configFilePath  string
	config          map[interface{}]interface{}
	tagGroupsByName map[string][]Tag
}

type Tag struct {
	tags.ITag
	defaultValue string
	filters      map[interface{}]interface{}
	matches      []interface{}
}

func (t *TagGroup) InitExternalTagGroups(configFilePath string) {
	t.configFilePath = configFilePath
	t.tagGroupsByName = make(map[string][]Tag)
	t.InitExternalTagGroup()

}

func (t *TagGroup) InitTagGroup(_ string, skippedTags []string) {
	t.SkippedTags = skippedTags
}

func (t *TagGroup) InitExternalTagGroup() {
	configMap := make(map[interface{}]interface{})
	confBytes, err := ioutil.ReadFile(t.configFilePath)
	err = yaml.Unmarshal(confBytes, &configMap)
	if err != nil {
		logger.Error(err.Error())
	}
	t.config = configMap
	t.extractExternalTags()
	//TODO
	//t.SetTags(externalTags)
}

func (t *TagGroup) extractExternalTags() {
	tagGroups := t.config["tag_group"]
	switch tgs := tagGroups.(type) {
	case map[interface{}]interface{}:
		groupName := tgs["name"].(string)
		t.tagGroupsByName[groupName] = t.ExtractExternalGroupsTags(tgs["tags"].([]interface{}))
	case []interface{}:
		for _, tagGroup := range tgs {
			tagGroupMap := tagGroup.(map[interface{}]interface{})
			groupName := tagGroupMap["name"].(string)
			t.tagGroupsByName[groupName] = t.ExtractExternalGroupsTags(tagGroupMap["tags"].([]interface{}))
		}
	}
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{}
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return t.UpdateBlockTags(block, struct{}{})
}

func (t *TagGroup) ExtractExternalGroupsTags(rawTags []interface{}) []Tag {
	var groupTags []Tag
	for _, rawTag := range rawTags {
		rawTagMap := rawTag.(map[interface{}]interface{})
		tagKey := rawTagMap["name"].(string)
		groupFilters := rawTagMap["filters"].(map[interface{}]interface{})
		tagValueObj := rawTagMap["value"]
		computedTag, err := calculateExternalTag(tagValueObj.(map[interface{}]interface{}), tagKey, groupFilters)
		if err != nil {
			logger.Error(err.Error())
		}
		groupTags = append(groupTags, computedTag)
	}
	return groupTags
}

func calculateExternalTag(tagValueObj map[interface{}]interface{}, tagKey string, groupFilters map[interface{}]interface{}) (Tag, error) {
	var calculatedTag = Tag{filters: groupFilters}
	var okDefault, okMatches bool
	var defaultValue, matches interface{}
	if defaultValue, okDefault = tagValueObj["default"]; okDefault {
		calculatedTag.defaultValue = defaultValue.(string)
		calculatedTag.ITag = &tags.Tag{Key: tagKey, Value: defaultValue.(string)}
	}
	if matches, okMatches = tagValueObj["matches"]; okMatches {
		calculatedTag.matches = matches.([]interface{})
	}
	if !okDefault && !okMatches {
		return Tag{}, errors.New("please specify either a default tag value and/or a computed tag value")
	}
	return calculatedTag, nil
}
