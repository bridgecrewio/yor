package simple

import (
	"os"
	"sort"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTagGroup(t *testing.T) {
	t.Run("test tagGroup CreateTagsForBlock", func(t *testing.T) {
		path := "../../../../tests/utils/blameutils/git_tagger_file.txt"
		tagGroup := TagGroup{}
		tagGroup.InitTagGroup("", nil, nil)

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
		tagGroup.SetTags(extraTags)
		block := &MockTestBlock{
			Block: structure.Block{
				FilePath:   path,
				IsTaggable: true,
			},
		}

		_ = tagGroup.CreateTagsForBlock(block)
		assert.Equal(t, len(block.NewTags), len(extraTags))
	})

	t.Run("Test create tags from env", func(t *testing.T) {
		tagGroup := TagGroup{}
		_ = os.Setenv("YOR_SIMPLE_TAGS", "{\"foo\": \"bar\", \"foo2\": \"bar2\"}")
		tagGroup.InitTagGroup("", nil, nil)
		getTags := tagGroup.GetTags()

		expected := []tags.Tag{{Key: "foo", Value: "bar"}, {Key: "foo2", Value: "bar2"}}
		sort.Slice(getTags, func(i, j int) bool {
			return getTags[i].GetKey() < getTags[j].GetKey()
		})

		assert.Equal(t, 2, len(getTags))
		for i, expectedTag := range expected {
			assert.Equal(t, expectedTag.Key, getTags[i].GetKey())
			assert.Equal(t, expectedTag.Value, getTags[i].GetValue())
		}
	})
}

type MockTestBlock struct {
	structure.Block
}

func (b *MockTestBlock) UpdateTags() {
}

func (b *MockTestBlock) Init(_ string, _ interface{}) {}

func (b *MockTestBlock) String() string {
	return ""
}

func (b *MockTestBlock) GetResourceID() string {
	return ""
}

func (b *MockTestBlock) GetLines(_ ...bool) structure.Lines {
	return structure.Lines{Start: 1, End: 3}
}

func (b *MockTestBlock) GetTagsLines() structure.Lines {
	return structure.Lines{Start: -1, End: -1}
}

func (b *MockTestBlock) GetSeparator() string {
	return ""
}
