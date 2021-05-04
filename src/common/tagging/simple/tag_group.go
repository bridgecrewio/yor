package simple

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/utils"
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

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) {
	utils.CreateTagsForBlock(t, block)
}
