package simple

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type TagGroup struct {
	tagging.TagGroup
}

func (t *TagGroup) InitTagGroup(_ string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
	t.SkippedTags = skippedTags
	t.SpecifiedTags = explicitlySpecifiedTags
	envTagsStr := os.Getenv("YOR_SIMPLE_TAGS")
	if envTagsStr == "" {
		return
	}
	logger.Debug(fmt.Sprintf("Simple tags from env: %v", envTagsStr))
	var extraTagsFromArgs map[string]string
	if strings.HasPrefix(envTagsStr, "'") {
		envTagsStr = envTagsStr[1 : len(envTagsStr)-1]
	}
	if strings.HasPrefix(envTagsStr, "\"") {
		err := json.Unmarshal([]byte(envTagsStr), &envTagsStr)
		if err != nil {
			logger.Info(fmt.Sprintf("failed to parse extra tags from env: %s", err))
		}
	}
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
func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return t.UpdateBlockTags(block, struct{}{})
}
