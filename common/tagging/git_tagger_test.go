package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	commonStructure "bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"bridgecrewio/yor/tests/utils/blameutils"
	"bridgecrewio/yor/tests/utils/structureutils"
	"testing"

	"github.com/go-git/go-git/v5"

	"github.com/stretchr/testify/assert"
)

func TestGitTagger(t *testing.T) {
	path := "test_file"
	blame := blameutils.SetupBlameResults(t, path)

	t.Run("test git tagger CreateTagsForBlock", func(t *testing.T) {
		gitService := &gitservice.GitService{
			BlameByFile: map[string]*git.BlameResult{path: blame},
		}
		tagger := GitTagger{Tagger: Tagger{
			Tags: []tags.ITag{},
		},
			GitService: gitService,
		}

		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
		}
		tagger.InitTags(extraTags)
		block := &structureutils.MockTestBlock{
			Block: commonStructure.Block{
				FilePath:   path,
				IsTaggable: true,
			},
		}

		tagger.CreateTagsForBlock(block)
		assert.Equal(t, len(block.NewTags), len(tags.TagTypes)+len(extraTags))
	})
}
