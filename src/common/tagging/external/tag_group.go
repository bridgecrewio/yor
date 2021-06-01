package external

import (
	"errors"
	"fmt"
	"io/ioutil"

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

func (t Tag) SatisfyFilters(block structure.IBlock) bool {
	newTags, existingTags := block.GetNewTags(), block.GetExistingTags()
	blockTags := append(newTags, existingTags...)
	satisfyFilters := true
	for filterKey, filterValue := range t.filters {
		if filterKey == "tags" {
			for filterTagKey, filterTagValue := range filterValue.(map[interface{}]interface{}) {
				foundFilterTag := false
				for _, blockTag := range blockTags {
					if blockTag.GetKey() == filterTagKey && blockTag.GetValue() == filterTagValue {
						foundFilterTag = true
						break
					}
				}
				satisfyFilters = satisfyFilters && foundFilterTag
			}
		}
		if filterKey == "directory" {
			// TODO Filter dir
			fmt.Println()
		}
	}
	return satisfyFilters
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
	blockTags := make([]tags.ITag, len(newTags)+len(existingTags))
	for _, groupTags := range t.tagGroupsByName {
		for _, groupTag := range groupTags {
			tagValue, err := calculateTagValue(block, groupTag)
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
	blockTags = append(filteredNewTags, existingTags...)
	t.SetTags(blockTags)
	block.AddNewTags(filteredNewTags)
	return nil
}

func calculateTagValue(block structure.IBlock, tag Tag) (tags.ITag, error) {
	var retTag = &tags.Tag{}
	if !tag.SatisfyFilters(block) {
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
	return Tag{}, errors.New(fmt.Sprintf("Could not compute external tag %s", tag.GetKey()))
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
