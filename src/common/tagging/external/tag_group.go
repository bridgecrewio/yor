package external

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/utils"
	"gopkg.in/yaml.v2"
)

var EnvVariableRegex = regexp.MustCompile(`\${env:([^\s]+)}`)

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
	TagGroups []struct {
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

type MatchesConfig []map[string]interface{}

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
	tagGroups := t.config.TagGroups
	for _, tagGroup := range tagGroups {
		logger.Info(fmt.Sprintf("extracting tag group named %v from yaml", tagGroup))
		tagGroupTags := tagGroup.Tags
		tagGroupName := evaluateTemplateVariable(tagGroup.TagGroupName)
		t.tagGroupsByName[tagGroupName] = t.ExtractExternalGroupsTags(tagGroupTags)
	}
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{}
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	logger.Info(fmt.Sprintf("extrnal tag group creating tags for block %v", block.GetResourceID()))
	newTags, existingTags := block.GetNewTags(), block.GetExistingTags()
	var filteredNewTags = make([]tags.ITag, len(newTags))
	blockTags := make([]tags.ITag, len(newTags)+len(existingTags))
	copy(filteredNewTags, newTags)
	newTagsNum := 0
	for _, groupTags := range t.tagGroupsByName {
		for _, groupTag := range groupTags {
			tagValue, err := t.CalculateTagValue(block, groupTag)
			if err != nil {
				logger.Error(err.Error())
			}
			newTagsNum++
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
	logger.Info(fmt.Sprintf("Created %d new tags", newTagsNum))
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
	retTag.Value = evaluateTemplateVariable(tag.defaultValue)
	blockTags := append(block.GetExistingTags(), block.GetNewTags()...)
	if len(tag.matches) > 0 {
		for _, matchEntry := range tag.matches {
			for matchValue, matchObj := range matchEntry {
				// Currently, we only allow matches on tags
				switch matchType := matchObj.(type) {
				case string:
					retTag.Value = evaluateTemplateVariable(matchType)
				case map[interface{}]interface{}:
					for tagName, tagMatch := range matchType["tags"].(map[interface{}]interface{}) {
						switch tagMatch.(type) {
						case string:
							for _, blockTag := range blockTags {
								blockTagKey, blockTagValue := blockTag.GetKey(), blockTag.GetValue()
								if blockTagKey == tagName && blockTagValue == tagMatch {
									retTag.Value = evaluateTemplateVariable(matchValue)
								}
							}
						case []interface{}:
							for _, blockTag := range blockTags {
								blockTagKey, blockTagValue := blockTag.GetKey(), blockTag.GetValue()
								if blockTagKey == tagName {
									if blockTagKey == tags.GitModifiersTagKey {
										for _, val := range strings.Split(blockTagValue, "/") {
											if utils.InSlice(tagMatch, val) {
												retTag.Value = matchValue
												break
											}
										}
									} else if utils.InSlice(tagMatch, blockTagValue) {
										retTag.Value = matchValue
									}
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
	return Tag{}, fmt.Errorf("could not compute external tag %s", tag.GetKey())
}

func (t *TagGroup) ExtractExternalGroupsTags(tagsConfig TagsConfig) []Tag {
	var groupTags []Tag
	for _, tagConfig := range tagsConfig {
		var groupFilters = tagConfig.Filters
		tagValueObj := tagConfig.TagValue
		tagKey := evaluateTemplateVariable(tagConfig.TagKey)
		computedTag, err := parseExternalTag(tagValueObj, tagKey, groupFilters)
		if err != nil {
			logger.Error(err.Error())
		}
		groupTags = append(groupTags, computedTag)
	}
	return groupTags
}

func evaluateTemplateVariable(val string) string {
	envVariableMatch := EnvVariableRegex.FindStringSubmatch(val)
	if len(envVariableMatch) == 2 {
		envVal, exists := os.LookupEnv(envVariableMatch[1])
		if !exists {
			logger.Warning(fmt.Sprintf("environment variable %s is not found", envVariableMatch[1]))
		} else {
			return envVal
		}
	}
	return val
}

func parseExternalTag(tagValueObj TagConfigValue, tagKey string, groupFilters FiltersConfig) (Tag, error) {
	var parsedTag = Tag{filters: groupFilters}
	if tagValueObj.Matches == nil && tagValueObj.Default == "" {
		return Tag{}, fmt.Errorf("please specify either a default tag value and/or a computed tag value")
	}
	parsedTag.defaultValue = tagValueObj.Default
	parsedTag.ITag = &tags.Tag{Key: tagKey, Value: tagValueObj.Default}
	parsedTag.matches = tagValueObj.Matches
	return parsedTag, nil
}
