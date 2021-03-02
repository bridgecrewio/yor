package simple

import (
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
)

type Tagger struct {
	tagging.Tagger
}

func (t *Tagger) InitTagger(_ string, skippedTags []string) {
	t.SkippedTags = skippedTags
}

func (t *Tagger) CreateTagsForBlock(block structure.IBlock) {
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
