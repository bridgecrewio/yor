package external

import (
	"os"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/spf13/viper"
)

type TagGroup struct {
	tagging.TagGroup
	configFilePath string
	config         map[string]interface{}
}

func (t *TagGroup) InitConfigFile(configFilePath string) {
	t.configFilePath = configFilePath
}

func (t *TagGroup) InitTagGroup(_ string, skippedTags []string) {
	t.SkippedTags = skippedTags
}

func (t *TagGroup) InitExternalTagGroup() {
	viper.SetConfigName("external_tags")
	viper.SetConfigType("yaml")
	file, err := os.Open(t.configFilePath)
	if err != nil {
		logger.Error(err.Error())
	}
	err = viper.ReadConfig(file)
	if err != nil {
		logger.Error(err.Error())
	}
	t.config = viper.AllSettings()
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return nil
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return t.UpdateBlockTags(block, struct{}{})
}
