package structure

import (
	"encoding/json"
	"fmt"
	goformationTags "github.com/awslabs/goformation/v5/cloudformation/tags"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
)

type CloudformationBlock struct {
	structure.Block
}

func (b *CloudformationBlock) UpdateTags() {
	if !b.IsTaggable {
		return
	}

	mergedTags := b.MergeTags()
	cfnMergedTags := make([]goformationTags.Tag, 0)
	for _, t := range mergedTags {
		cfnMergedTags = append(cfnMergedTags, goformationTags.Tag{
			Key:   t.GetKey(),
			Value: t.GetValue(),
		})
	}

	blockBytes, _ := json.Marshal(b.RawBlock)
	var blockAsMap map[string]interface{}
	err := json.Unmarshal(blockBytes, &blockAsMap)
	if err != nil {
		logger.Warning(fmt.Sprintf("failed to marshal block to json: %s", err))
		return
	}

	blockAsMap["Properties"].(map[string]interface{})["Tags"] = cfnMergedTags
	b.RawBlock = blockAsMap
}

func (b *CloudformationBlock) GetTagsLines() structure.Lines {
	return b.TagLines
}

func (b *CloudformationBlock) GetSeparator() string {
	return "/n"
}
