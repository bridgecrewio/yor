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
			&tags.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
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
		}

		diff := block.CalculateTagsDiff()
		merged := block.MergeTags()

		assert.Equal(t, 3, len(merged), "Merging failed, expected to see 3 tags")
		assert.Equal(t, newTags[0].GetValue(), diff.Updated[0].NewValue)
		assert.Equal(t, newTags[1].GetValue(), diff.Added[0].GetValue())
	})
	t.Run("Test no reported diff for non-yor tags diff", func(t *testing.T) {
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
			&tags.Tag{
				Key:   "env",
				Value: "dev",
			},
			&tags.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
		}

		newTags := []tags.ITag{
			&tags.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf",
				},
			},
			&tags.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
			&tags.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
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
		}

		diff := block.CalculateTagsDiff()
		merged := block.MergeTags()

		assert.Equal(t, 3, len(merged), "Merging failed, expected to see 3 tags")
		assert.Equal(t, 0, len(diff.Updated))
		assert.Equal(t, 0, len(diff.Added))
	})

	t.Run("Ensure old trace tag is not overridden by a new trace tag", func(t *testing.T) {
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
			&tags.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
			&tags.YorTraceTag{
				Tag: tags.Tag{
					Key:   "yor_trace",
					Value: "my-old-trace",
				},
			},
		}
		newTags := []tags.ITag{
			&tags.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "hatulik",
				},
			},
			&tags.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
			&tags.YorTraceTag{
				Tag: tags.Tag{
					Key:   "yor_trace",
					Value: "my-new-trace",
				},
			},
		}

		block := structure.TerraformBlock{
			Block: structure2.Block{
				FilePath:          "",
				ExitingTags:       existingTags,
				NewTags:           []tags.ITag{},
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "",
			},
		}

		block.AddNewTags(newTags)
		diff := block.CalculateTagsDiff()
		merged := block.MergeTags()
		assert.Equal(t, 1, len(diff.Updated))
		for _, tag := range merged {
			if traceTag, ok := tag.(*tags.YorTraceTag); ok {
				assert.Equal(t, traceTag.Value, "my-old-trace")
			}
		}
	})

}
