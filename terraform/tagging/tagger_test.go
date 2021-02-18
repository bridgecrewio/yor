package tagging

import (
	commonStructure "bridgecrewio/yor/common/structure"
	commonTagging "bridgecrewio/yor/common/tagging"
	"bridgecrewio/yor/common/tagging/tags"
	"bridgecrewio/yor/terraform/structure"
	"bridgecrewio/yor/tests/utils/blameUtils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformTagger(t *testing.T) {
	blame := blameUtils.SetupBlame(t)

	t.Run("test terraform tagger CreateTagsForBlock", func(t *testing.T) {
		tagger := TerraformTagger{Tagger: commonTagging.Tagger{
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
			Block: commonStructure.Block{
				FilePath:          "",
				ExitingTags:       []tags.ITag{},
				NewTags:           []tags.ITag{},
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "",
			},
		}
		tagger.CreateTagsForBlock(block, &blame)
		assert.Equal(t, len(block.NewTags), len(tags.TagTypes)+len(extraTags))
	})
}
