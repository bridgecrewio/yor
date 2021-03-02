package code2cloud

import (
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
)

type Tagger struct {
	tagging.Tagger
	traceTag YorTraceTag
}

func (t *Tagger) InitTagger(_ string, skippedTags []string) {
	t.SkippedTags = skippedTags
	t.SetTags([]tags.ITag{&YorTraceTag{}})
}

func (t *Tagger) CreateTagsForBlock(block structure.IBlock) {
	var newTags []tags.ITag
	for _, tag := range t.GetTags() {
		tagVal, err := tag.CalculateValue(struct{}{})
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create %v tag for block %v", tag.GetKey(), block.GetResourceID()))
		}
		newTags = append(newTags, tagVal)
	}
	block.AddNewTags(newTags)
}
