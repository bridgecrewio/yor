package simple

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
	"encoding/json"
	"fmt"
	"os"
)

type TagGroup struct {
	tagging.TagGroup
}

func (t *TagGroup) InitTagGroup(_ string, skippedTags []string) {
	t.SkippedTags = skippedTags
	envTagsStr := os.Getenv("YOR_SIMPLE_TAGS")
	var extraTagsFromArgs map[string]string
	if err := json.Unmarshal([]byte(envTagsStr), &extraTagsFromArgs); err != nil {
		logger.Info(fmt.Sprintf("failed to parse extra tags from env: %s", err))
	} else {
		var envTags []tags.ITag
		for key, value := range extraTagsFromArgs {
			envTags = append(envTags, tags.Init(key, value))
		}

		t.SetTags(envTags)
	}
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{}
}

func (t *TagGroup) CreateTagsForBlock(block common.IBlock) {
	var newTags []tags.ITag
	for _, tag := range t.GetTags() {
		tagVal, err := tag.CalculateValue(struct{}{})
		if err != nil {
			logger.Warning("Failed to create extra tag", tag.GetKey())
		}
		newTags = append(newTags, tagVal)
	}
	block.AddNewTags(newTags)
}
