package external

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/utils"
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

func (t Tag) SatisfyFilters(block structure.IBlock, tagFilterDir string) bool {
	newTags, existingTags := block.GetNewTags(), block.GetExistingTags()
	newTags = append(newTags, existingTags...)
	satisfyFilters := true
	for filterKey, filterValue := range t.filters {
		if filterKey == "tags" {
			for filterTagKey, filterTagValue := range filterValue.(map[interface{}]interface{}) {
				strFilterValue := filterTagValue
				if val, ok := filterTagValue.(int); ok {
					strFilterValue = strconv.Itoa(val)
				}
				foundFilterTag := false
				for _, blockTag := range newTags {
					if blockTag.GetKey() == filterTagKey && blockTag.GetValue() == strFilterValue {
						foundFilterTag = true
						break
					}
				}
				satisfyFilters = satisfyFilters && foundFilterTag
			}
		}
		if filterKey == "directory" {
			if tagFilterDir != filterValue {
				satisfyFilters = false
				break
			}
		}
	}
	return satisfyFilters
}

func (t *TagGroup) InitExternalTagGroups(configFilePath string) {
	t.configFilePath = configFilePath
	t.tagGroupsByName = make(map[string][]Tag)
	t.InitExternalTagGroup()

}

func (t *TagGroup) InitTagGroup(dir string, skippedTags []string) {
	t.SkippedTags = skippedTags
	t.Dir = dir
}

func (t *TagGroup) InitExternalTagGroup() {
	configMap := make(map[interface{}]interface{})
	confBytes, err := ioutil.ReadFile(t.configFilePath)
	if err != nil {
		logger.Error(err.Error())
	}
	errYaml := yaml.Unmarshal(confBytes, &configMap)
	if errYaml != nil {
		logger.Error(errYaml.Error())
	}
	t.config = configMap
	t.extractExternalTags()
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
	newTags, existingTags := block.GetNewTags(), block.GetExistingTags()
	var filteredNewTags = make([]tags.ITag, len(newTags))
	copy(filteredNewTags, newTags)
	for _, groupTags := range t.tagGroupsByName {
		for _, groupTag := range groupTags {
			tagValue, err := t.calculateTagValue(block, groupTag)
			if err != nil {
				logger.Error(err.Error())
			}
			if tagValue == nil {
				for i, newTag := range newTags {
					if newTag.GetKey() == groupTag.GetKey() {
						filteredNewTags = append(filteredNewTags[:i], filteredNewTags[i+1:]...)
					}
				}
			} else {
				filteredNewTags = append(filteredNewTags, tagValue)
			}
		}
	}
	filteredNewTags = append(filteredNewTags, existingTags...)
	t.SetTags(filteredNewTags)
	block.AddNewTags(filteredNewTags)
	return nil
}

func (t *TagGroup) calculateTagValue(block structure.IBlock, tag Tag) (tags.ITag, error) {
	var retTag = &tags.Tag{}
	if !tag.SatisfyFilters(block, t.Dir) {
		return nil, nil
	}
	retTag.Key = tag.GetKey()
	retTag.Value = tag.defaultValue
	blockTags := append(block.GetExistingTags(), block.GetNewTags()...)
	if len(tag.matches) > 0 {
		for _, matchEntry := range tag.matches {
			for matchValue, matchObj := range matchEntry.(map[interface{}]interface{}) {
				// Currently, we only allow matches on tags
				if matchTags, ok := matchObj.(map[interface{}]interface{})["tags"]; ok {
					for tagName, tagMatch := range matchTags.(map[interface{}]interface{}) {
						switch match := tagMatch.(type) {
						case string:
							for _, blockTag := range blockTags {
								blockTagKey, blockTagValue := blockTag.GetKey(), blockTag.GetValue()
								if blockTagKey == tagName && blockTagValue == match {
									retTag.Value = matchValue.(string)
								}
							}
						case []interface{}:
							for _, blockTag := range blockTags {
								blockTagKey, blockTagValue := blockTag.GetKey(), blockTag.GetValue()
								if blockTagKey == tagName && utils.InSlice(match, blockTagValue) {
									retTag.Value = matchValue.(string)
								}
							}
						}
					}
				}
			}
		}
		return retTag, nil
	} else if tag.defaultValue != "" {
		return retTag, nil
	}
	return Tag{}, fmt.Errorf("Could not compute external tag %s", tag.GetKey())
}

func (t *TagGroup) ExtractExternalGroupsTags(rawTags []interface{}) []Tag {
	var groupTags []Tag
	for _, rawTag := range rawTags {
		var groupFilters map[interface{}]interface{}
		rawTagMap := rawTag.(map[interface{}]interface{})
		tagKey := rawTagMap["name"].(string)
		if _, okFilters := rawTagMap["filters"]; okFilters {
			groupFilters = rawTagMap["filters"].(map[interface{}]interface{})
		}
		tagValueObj := rawTagMap["value"]
		computedTag, err := parseExternalTag(tagValueObj.(map[interface{}]interface{}), tagKey, groupFilters)
		if err != nil {
			logger.Error(err.Error())
		}
		groupTags = append(groupTags, computedTag)
	}
	return groupTags
}

func parseExternalTag(tagValueObj map[interface{}]interface{}, tagKey string, groupFilters map[interface{}]interface{}) (Tag, error) {
	var parsedTag = Tag{filters: groupFilters}
	var okDefault, okMatches bool
	var defaultValue, matches interface{}
	if defaultValue, okDefault = tagValueObj["default"]; okDefault {
		parsedTag.defaultValue = defaultValue.(string)
		parsedTag.ITag = &tags.Tag{Key: tagKey, Value: defaultValue.(string)}
	}
	if matches, okMatches = tagValueObj["matches"]; okMatches {
		parsedTag.matches = matches.([]interface{})
	}
	if !okDefault && !okMatches {
		return Tag{}, errors.New("please specify either a default tag value and/or a computed tag value")
	}
	return parsedTag, nil
}
