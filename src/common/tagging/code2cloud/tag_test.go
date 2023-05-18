package code2cloud

import (
	"regexp"
	"strings"
	"testing"

	"github.com/bridgecrewio/yor/src/common/tagging"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"

	"github.com/stretchr/testify/assert"
)

func TestTagCreation(t *testing.T) {
	t.Run("BcTraceTagCreation", func(t *testing.T) {
		tag := YorTraceTag{}
		valueTag := EvaluateTag(t, &tag)
		assert.Equal(t, "yor_trace", valueTag.GetKey())
		assert.Equal(t, 36, len(valueTag.GetValue()))
	})
	t.Run("BcTraceTagCreationWithPrefix", func(t *testing.T) {
		tag := YorTraceTag{}
		tagPrefix := "prefix_"
		valueTag := EvaluateTagWithPrefix(t, &tag, tagPrefix)
		assert.Equal(t, tagPrefix+"yor_trace", valueTag.GetKey())
		assert.Equal(t, 36, len(valueTag.GetValue()))
	})
}

func EvaluateTag(t *testing.T, tag tags.ITag) tags.ITag {
	tag.Init()
	newTag, err := tag.CalculateValue(struct{}{})
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &tags.Tag{}, newTag)
	assert.True(t, IsValidUUID(newTag.GetValue()))

	return newTag
}

func EvaluateTagWithPrefix(t *testing.T, tag tags.ITag, tagPrefix string) tags.ITag {
	tag.Init()
	tag.SetTagPrefix(tagPrefix)
	newTag, err := tag.CalculateValue(struct{}{})
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &tags.Tag{}, newTag)
	assert.True(t, IsValidUUID(newTag.GetValue()))

	return newTag
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func TestCode2CloudTagGroup(t *testing.T) {
	t.Run("test tagGroup CreateTagsForBlock", func(t *testing.T) {
		path := "../../../../tests/utils/blameutils/git_tagger_file.txt"
		tagGroup := TagGroup{}
		tagGroup.InitTagGroup("", nil, nil, tagging.WithTagPrefix("prefix"))

		blockName := "my_bucket"

		block := &MockTestBlock{
			Block: structure.Block{
				FilePath:   path,
				IsTaggable: true,
				Name:       blockName,
			},
		}

		_ = tagGroup.CreateTagsForBlock(block)
		assert.Equal(t, 2, len(block.NewTags))
		foundYorName := false
		for _, newTag := range block.NewTags {
			if strings.Contains(newTag.GetKey(), tags.YorNameTagKey) {
				foundYorName = true
				assert.Equal(t, blockName, newTag.GetValue())
			}
		}
		assert.True(t, foundYorName)
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
