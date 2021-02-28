package simple

import (
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
)

type Tagger struct {
	tagging.Tagger
	extraTags []tags.ITag
}

func (t *Tagger) InitTagger(_ string) {
}

func (t *Tagger) InitExtraTags(extraTags []tags.ITag) {
	t.extraTags = append(t.extraTags, extraTags...)
	for _, val := range t.extraTags {
		val.Init()
	}
}

func (t *Tagger) CreateTagsForBlock(block structure.IBlock) {
	var newTags []tags.ITag
	for _, tag := range t.extraTags {
		tagVal, err := tag.CalculateValue(struct{}{})
		if err != nil {
			logger.Warning("Failed to create extra tag", tag.GetKey())
		}
		newTags = append(newTags, tagVal)
	}
	block.AddNewTags(newTags)
}
