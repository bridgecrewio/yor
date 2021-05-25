package external

import (
	"io/ioutil"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"gopkg.in/yaml.v2"
)

type TagGroup struct {
	tagging.TagGroup
	configFilePath string
	config         map[interface{}]interface{}
}

func (t *TagGroup) InitConfigFile(configFilePath string) {
	t.configFilePath = configFilePath
}

func (t *TagGroup) InitTagGroup(_ string, skippedTags []string) {
	t.SkippedTags = skippedTags
}

func (t *TagGroup) InitExternalTagGroup() {
	m := make(map[interface{}]interface{})
	confBytes, err := ioutil.ReadFile(t.configFilePath)
	err = yaml.Unmarshal(confBytes, &m)
	if err != nil {
		logger.Error(err.Error())
	}
	t.config = m
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return nil
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return t.UpdateBlockTags(block, struct{}{})
}
