package code2cloud

import (
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
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
	if len(t.GetTags()) > 0 {
		tag, err := t.GetTags()[0].CalculateValue(struct{}{})
		if err != nil {
			logger.Error("Failed to create yor trace tag for block", block.GetResourceID())
		}
		block.AddNewTags([]tags.ITag{tag})
	}
}
