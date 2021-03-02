package simple

import (
	"bridgecrewio/yor/src/common"
	commonStructure "bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTagger(t *testing.T) {
	t.Run("test tagger CreateTagsForBlock", func(t *testing.T) {
		path := "../../../../tests/utils/blameutils/git_tagger_file.txt"
		tagger := Tagger{}
		tagger.InitTagger("", nil)

		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
			&tags.Tag{
				Key:   "another_custom_value",
				Value: "custom",
			},
		}
		tagger.SetTags(extraTags)
		block := &MockTestBlock{
			Block: commonStructure.Block{
				FilePath:   path,
				IsTaggable: true,
			},
		}

		tagger.CreateTagsForBlock(block)
		assert.Equal(t, len(block.NewTags), len(extraTags))
	})
}

type MockTestBlock struct {
	commonStructure.Block
}

func (b *MockTestBlock) Init(_ string, _ interface{}) {}

func (b *MockTestBlock) String() string {
	return ""
}

func (b *MockTestBlock) GetResourceID() string {
	return ""
}

func (b *MockTestBlock) GetLines(_ ...bool) common.Lines {
	return common.Lines{Start: 1, End: 3}
}
