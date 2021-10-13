package external

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/stretchr/testify/assert"
)

func TestExternalTagGroup(t *testing.T) {

	t.Run("test tagGroup CreateTagsForBlock default value", func(t *testing.T) {
		_ = os.Setenv("GIT_BRANCH", "master")
		confPath, _ := filepath.Abs("../../../../tests/external_tags/external_tag_group.yml")
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
						Value: "checkov",
					},
				},
			},
		}
		err := tagGroup.CreateTagsForBlock(block)
		if err != nil {
			logger.Warning(err.Error())
			t.Fail()
		}
		for _, newBlockTag := range block.GetNewTags() {
			if newBlockTag.GetKey() == "env" {
				assert.Equal(t, "master", newBlockTag.GetValue())
			}
		}
		assert.Equal(t, 1, len(block.NewTags))
	})

	t.Run("test tagGroup CreateTagsForBlock matches", func(t *testing.T) {
		confPath, _ := filepath.Abs("../../../../tests/external_tags/external_tag_group.yml")
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
		assert.Equal(t, 1, len(block.NewTags))
		if len(block.NewTags) == 1 {
			assert.Equal(t, "env", block.NewTags[0].GetKey())
			assert.Equal(t, "dev", block.NewTags[0].GetValue())
		}
	})

	t.Run("test tagGroup CreateTagsForBlock matches with directory filter", func(t *testing.T) {
		confPath, _ := filepath.Abs("../../../../tests/external_tags/external_tag_group_dir.yml")
		tagGroup := TagGroup{}
		tagGroup.InitTagGroup("", nil)
		tagGroup.InitExternalTagGroups(confPath)
		block := &MockTestBlock{
			Block: structure.Block{
				FilePath:   "src/account/main.tf",
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
		assert.Equal(t, 2, len(block.NewTags))
		var dirTag tags.ITag
		for _, t := range block.NewTags {
			if t.GetKey() == "stack" {
				dirTag = t
				break
			}
		}
		assert.NotNil(t, dirTag)
		if dirTag != nil {
			assert.Equal(t, dirTag.GetKey(), "stack")
			assert.Equal(t, dirTag.GetValue(), "account")
		}
	})

	t.Run("test tagGroup CreateTagsForBlock not matches with directory filter", func(t *testing.T) {
		confPath, _ := filepath.Abs("../../../../tests/external_tags/external_tag_group_dir.yml")
		tagGroup := TagGroup{}
		tagGroup.InitTagGroup("", nil)
		tagGroup.InitExternalTagGroups(confPath)
		block := &MockTestBlock{
			Block: structure.Block{
				FilePath:   "src/base/main.tf",
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
		assert.Equal(t, 4, len(block.ExitingTags))
		assert.Equal(t, 1, len(block.NewTags))
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
