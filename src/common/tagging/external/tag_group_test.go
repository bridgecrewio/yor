package external

import (
	"testing"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/stretchr/testify/assert"
)

func TestSimpleTagGroup(t *testing.T) {

	t.Run("test tagGroup CreateTagsForBlock default value", func(t *testing.T) {
		confPath := "external_tag_group.yml"
		tagGroup := TagGroup{}
		tagGroup.InitTagGroup("", nil)
		tagGroup.InitExternalTagGroups(confPath)
		block := &MockTestBlock{
			Block: structure.Block{
				FilePath:   "",
				IsTaggable: true,
				ExitingTags: []tags.ITag{
					&tags.Tag{
						Key:   "git_modifiers",
						Value: "tronxd",
					},
					&tags.Tag{
						Key:   "git_repo",
						Value: "yor",
					},
				},
			},
		}
		err := tagGroup.CreateTagsForBlock(block)
		if err != nil {
			logger.Warning(err.Error())
			t.Fail()
		}
		assert.Equal(t, 3, len(block.ExitingTags)+len(block.NewTags))
	})

	t.Run("test tagGroup CreateTagsForBlock matches", func(t *testing.T) {
		confPath := "external_tag_group.yml"
		tagGroup := TagGroup{}
		tagGroup.InitTagGroup("", nil)
		tagGroup.InitExternalTagGroups(confPath)
		block := &MockTestBlock{
			Block: structure.Block{
				FilePath:   "",
				IsTaggable: true,
				ExitingTags: []tags.ITag{
					&tags.Tag{
						Key:   "git_modifiers",
						Value: "tronxd",
					},
					&tags.Tag{
						Key:   "git_repo",
						Value: "yor",
					},
					&tags.Tag{
						Key:   "git_commit",
						Value: "asd12f",
					},
					&tags.Tag{
						Key:   "yor_trace",
						Value: "123",
					},
				},
			},
		}
		err := tagGroup.CreateTagsForBlock(block)
		if err != nil {
			logger.Warning(err.Error())
			t.Fail()
		}
		assert.Equal(t, 5, len(block.ExitingTags)+len(block.NewTags))
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
