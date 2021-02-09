package structure

import (
	structure2 "bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"bridgecrewio/yor/terraform/structure"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerrraformBlock(t *testing.T) {
	t.Run("Test tag merging and diff", func(t *testing.T) {
		existingTags := []tags.ITag{
			&tags.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf",
				},
			},
			&tags.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
				},
			},
		}

		newTags := []tags.ITag{
			&tags.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf/hatulik",
				},
			},
			&tags.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
		}
		block := structure.TerraformBlock{
			Block: structure2.Block{
				FilePath:          "",
				ExitingTags:       existingTags,
				NewTags:           newTags,
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "",
			},
			NewOwner:      "",
			PreviousOwner: "",
			TraceId:       "",
		}

		diff := block.CalculateTagsDiff()
		merged := block.MergeTags()

		assert.Equal(t, 3, len(merged), "Merging failed, expected to see 3 tags")
		assert.Equal(t, newTags[0].GetValue(), diff["updated"][0].GetValue())
		assert.Equal(t, newTags[1].GetValue(), diff["added"][0].GetValue())
	})
}
