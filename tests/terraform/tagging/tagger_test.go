package tagging

import (
	structure2 "bridgecrewio/yor/common/structure"
	tagging2 "bridgecrewio/yor/common/tagging"
	"bridgecrewio/yor/common/tagging/tags"
	"bridgecrewio/yor/terraform/structure"
	"bridgecrewio/yor/terraform/tagging"
	"bridgecrewio/yor/tests/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerraformTagger(t *testing.T) {
	blame := utils.SetupBlame(t)

	t.Run("test terraform tagger IsBlockTaggable", func(t *testing.T) {
		tagger := tagging.TerraformTagger{}
		block := structure.TerraformBlock{
			Block: structure2.Block{
				FilePath:          "",
				ExitingTags:       []tags.ITag{},
				NewTags:           []tags.ITag{},
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "",
			},
			NewOwner:      "",
			PreviousOwner: "",
			TraceId:       "",
		}
		assert.True(t, tagger.IsBlockTaggable(&block))
	})

	t.Run("test terraform tagger CreateTagsForBlock", func(t *testing.T) {
		tagger := tagging.TerraformTagger{Tagger: tagging2.Tagger{
			Tags: []tags.ITag{},
		}}
		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
		}
		tagger.InitTags(extraTags)
		block := &structure.TerraformBlock{
			Block: structure2.Block{
				FilePath:          "",
				ExitingTags:       []tags.ITag{},
				NewTags:           []tags.ITag{},
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "",
			},
			NewOwner:      "",
			PreviousOwner: "",
			TraceId:       "",
		}
		err := tagger.CreateTagsForBlock(block, &blame)
		if err != nil {
			assert.Fail(t, "Failed to create tags for block", err)
		}
		assert.Equal(t, len(block.NewTags), len(tags.TagTypes)+len(extraTags))

	})
}
