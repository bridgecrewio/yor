package external

import (
	"fmt"
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
	tagGroupsByName map[string][]tags.ITag
}

func (t *TagGroup) InitExternalTagGroups(configFilePath string) {
	t.configFilePath = configFilePath
	t.tagGroupsByName = make(map[string][]tags.ITag)
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
	externalTags := t.extractExternalTags()
	t.SetTags(externalTags)
}

func (t *TagGroup) extractExternalTags() []tags.ITag {
	externalGroupTags := make([]tags.ITag, 0)
	tagGroups := t.config["tag_group"]
	switch tgs := tagGroups.(type) {
	case map[interface{}]interface{}:
		externalGroupTags = t.ExtractExternalGroupTags(tgs["tags"].([]interface{}))
	case []interface{}:
		for _, tagGroup := range tgs {
			fmt.Println(tagGroup)
			//externalGroupTags = append(externalGroupTags, t.extractExternalTags(tg.([]interface{}))...)
		}
	}
	return externalGroupTags
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{}
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return t.UpdateBlockTags(block, struct{}{})
}

func (t *TagGroup) ExtractExternalGroupTags(rawTags []interface{}) []tags.ITag {
	var tagGroups tags.ITag
	fmt.Println(tagGroups)
	return nil
	//TODO
}
