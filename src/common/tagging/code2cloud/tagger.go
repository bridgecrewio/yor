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

func (t *Tagger) InitTagger(_ string) {
	t.traceTag.Init()
}

func (t *Tagger) CreateTagsForBlock(block structure.IBlock) {
	tag, err := t.traceTag.CalculateValue(struct{}{})
	if err != nil {
		logger.Error("Failed to create yor trace tag for block", block.GetResourceID())
	}
	block.AddNewTags([]tags.ITag{tag})
}
