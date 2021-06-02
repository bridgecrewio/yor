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
	config          *Config
	tagGroupsByName map[string][]Tag
}

type Tag struct {
	tags.ITag
	defaultValue string
	filters      FiltersConfig
	matches      MatchesConfig
}

type Config struct {
	TagGroup []struct {
		TagGroupName string     `yaml:"name"`
		Tags         TagsConfig `yaml:"tags"`
	} `yaml:"tag_groups"`
}

type TagsConfig []struct {
	TagKey   string         `yaml:"name"`
	TagValue TagConfigValue `yaml:"value"`
	Filters  FiltersConfig  `yaml:"filters"`
}

type TagConfigValue struct {
	Default string        `yaml:"default"`
	Matches MatchesConfig `yaml:"matches"`
}

type MatchesConfig []interface{}

type FiltersConfig struct{ Tags map[string]interface{} }

func (t Tag) SatisfyFilters(block structure.IBlock, tagFilterDir string) bool {
	newTags, existingTags := block.GetNewTags(), block.GetExistingTags()
	var blockTags = make([]tags.ITag, len(newTags)+len(existingTags))
	copy(blockTags, append(newTags, existingTags...))
	satisfyFilters := true
	for filterKey, filterValue := range t.filters.Tags {
		switch filterKey {
		case "tags":
			for filterTagKey, filterTagValue := range filterValue.(map[interface{}]interface{}) {
				strFilterValue := filterTagValue
				if val, ok := filterTagValue.(int); ok {
					strFilterValue = strconv.Itoa(val)
				}
				foundFilterTag := false
				for _, blockTag := range blockTags {
					if blockTag.GetKey() == filterTagKey && blockTag.GetValue() == strFilterValue {
						foundFilterTag = true
						break
					}
				}
				satisfyFilters = satisfyFilters && foundFilterTag
			}

		case "directory":
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
	configMap := Config{}
	confBytes, err := ioutil.ReadFile(t.configFilePath)
	if err != nil {
		logger.Error(err.Error())
	}
	errYaml := yaml.Unmarshal(confBytes, &configMap)
	if errYaml != nil {
		logger.Error(errYaml.Error())
	}
	t.config = &configMap
	t.extractExternalTags()
}

func (t *TagGroup) extractExternalTags() {
	tagGroups := t.config.TagGroup
	for _, tagGroup := range tagGroups {
		tagGroupTags := tagGroup.Tags
		tagGroupName := tagGroup.TagGroupName
		t.tagGroupsByName[tagGroupName] = t.ExtractExternalGroupsTags(tagGroupTags)
	}
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{}
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	newTags, existingTags := block.GetNewTags(), block.GetExistingTags()
	var filteredNewTags = make([]tags.ITag, len(newTags))
	blockTags := make([]tags.ITag, len(newTags)+len(existingTags))
	copy(filteredNewTags, newTags)
	for _, groupTags := range t.tagGroupsByName {
		for _, groupTag := range groupTags {
			tagValue, err := t.CalculateTagValue(block, groupTag)
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
	copy(blockTags, append(filteredNewTags, existingTags...))
	t.SetTags(blockTags)
	block.AddNewTags(filteredNewTags)
	return nil
}

func (t *TagGroup) CalculateTagValue(block structure.IBlock, tag Tag) (tags.ITag, error) {
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
				matchMap := matchObj.(map[interface{}]interface{})
				for tagName, tagMatch := range matchMap["tags"].(map[interface{}]interface{}) {
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
		return retTag, nil
	} else if tag.defaultValue != "" {
		return retTag, nil
	}
	return Tag{}, fmt.Errorf("Could not compute external tag %s", tag.GetKey())
}

func (t *TagGroup) ExtractExternalGroupsTags(tagsConfig TagsConfig) []Tag {
	var groupTags []Tag
	for _, tagConfig := range tagsConfig {
		var groupFilters = tagConfig.Filters
		tagValueObj := tagConfig.TagValue
		tagKey := tagConfig.TagKey
		computedTag, err := parseExternalTag(tagValueObj, tagKey, groupFilters)
		if err != nil {
			logger.Error(err.Error())
		}
		groupTags = append(groupTags, computedTag)
	}
	return groupTags
}

func parseExternalTag(tagValueObj TagConfigValue, tagKey string, groupFilters FiltersConfig) (Tag, error) {
	var parsedTag = Tag{filters: groupFilters}
	if tagValueObj.Matches == nil && tagValueObj.Default == "" {
		return Tag{}, errors.New("please specify either a default tag value and/or a computed tag value")
	}
	parsedTag.defaultValue = tagValueObj.Default
	parsedTag.ITag = &tags.Tag{Key: tagKey, Value: tagValueObj.Default}
	parsedTag.matches = tagValueObj.Matches

	return parsedTag, nil
}
