package structure

import (
	structure2 "bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/code2cloud"
	"bridgecrewio/yor/src/common/tagging/git"
	"bridgecrewio/yor/src/common/tagging/tags"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerrraformBlock(t *testing.T) {
	t.Run("Test tag merging and diff", func(t *testing.T) {
		existingTags := []tags.ITag{
			&git.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf",
				},
			},
			&git.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
				},
			},
		}

		newTags := []tags.ITag{
			&git.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf/hatulik",
				},
			},
			&git.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
			&git.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
				},
			},
		}
		block := TerraformBlock{
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
			&git.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf",
				},
			},
			&git.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
				},
			},
			&tags.Tag{
				Key:   "env",
				Value: "dev",
			},
			&git.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
		}

		newTags := []tags.ITag{
			&git.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf",
				},
			},
			&git.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
			&git.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
				},
			},
		}
		block := TerraformBlock{
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
			&git.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "gandalf",
				},
			},
			&git.GitOrgTag{
				Tag: tags.Tag{
					Key:   "git_org",
					Value: "bridgecrewio",
				},
			},
			&git.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
			&code2cloud.YorTraceTag{
				Tag: tags.Tag{
					Key:   "yor_trace",
					Value: "my-old-trace",
				},
			},
		}
		newTags := []tags.ITag{
			&git.GitModifiersTag{
				Tag: tags.Tag{
					Key:   "git_modifiers",
					Value: "hatulik",
				},
			},
			&git.GitRepoTag{
				Tag: tags.Tag{
					Key:   "git_repository",
					Value: "terragoat",
				},
			},
			&code2cloud.YorTraceTag{
				Tag: tags.Tag{
					Key:   "yor_trace",
					Value: "my-new-trace",
				},
			},
		}

		block := TerraformBlock{
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
			if traceTag, ok := tag.(*code2cloud.YorTraceTag); ok {
				assert.Equal(t, traceTag.Value, "my-old-trace")
			}
		}
	})

}
