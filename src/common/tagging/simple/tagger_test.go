package simple

import (
	commonStructure "bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTagger(t *testing.T) {
	t.Run("test tagger CreateTagsForBlock", func(t *testing.T) {
		path := "../../../../tests/utils/blameutils/git_tagger_file.txt"
		tagger := Tagger{
			Tagger: tagging.Tagger{
				Tags: []tags.ITag{},
			},
			extraTags: []tags.ITag{},
		}
		tagger.InitTagger("")

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
		tagger.InitExtraTags(extraTags)
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

func (b *MockTestBlock) GetLines() []int {
	return []int{1, 3}
}
