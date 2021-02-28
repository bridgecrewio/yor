package code2cloud

import (
	commonStructure "bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagCreation(t *testing.T) {
	t.Run("BcTraceTagCreation", func(t *testing.T) {
		tag := YorTraceTag{}
		valueTag := EvaluateTag(t, &tag)
		assert.Equal(t, "yor_trace", valueTag.GetKey())
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

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func TestCode2CloudTagger(t *testing.T) {
	t.Run("test tagger CreateTagsForBlock", func(t *testing.T) {
		path := "../../../../tests/utils/blameutils/git_tagger_file.txt"
		tagger := Tagger{
			Tagger: tagging.Tagger{
				Tags: []tags.ITag{},
			},
		}
		tagger.InitTagger("")

		block := &MockTestBlock{
			Block: commonStructure.Block{
				FilePath:   path,
				IsTaggable: true,
			},
		}

		tagger.CreateTagsForBlock(block)
		assert.Equal(t, 1, len(block.NewTags))
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
