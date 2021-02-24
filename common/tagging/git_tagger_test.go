package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	commonStructure "bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"bridgecrewio/yor/tests/terraform/resources/taggedkms"
	"bridgecrewio/yor/tests/utils/blameutils"
	"bridgecrewio/yor/tests/utils/structureutils"
	"testing"

	"github.com/go-git/go-git/v5"

	"github.com/stretchr/testify/assert"
)

func TestGitTagger(t *testing.T) {
	path := "../../tests/utils/blameutils/git_tagger_file.txt"
	blame := blameutils.SetupBlameResults(t, path, 3)

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

func TestGitTagger_mapOriginFileToGitFile(t *testing.T) {
	t.Run("map tagged kms", func(t *testing.T) {
		expectedMapping := taggedkms.ExpectedFileMappingTagged
		gitTagger := GitTagger{}
		filePath := "../../tests/terraform/resources/taggedkms/tagged_kms.tf"
		blame := taggedkms.KmsBlame
		gitTagger.mapOriginFileToGitFile(filePath, &blame)
		assert.Equal(t, expectedMapping["originToGit"], gitTagger.fileLinesMapper[filePath].originToGit)
		assert.Equal(t, expectedMapping["gitToOrigin"], gitTagger.fileLinesMapper[filePath].gitToOrigin)
	})
	t.Run("map kms with deleted lines", func(t *testing.T) {
		expectedMapping := taggedkms.ExpectedFileMappingDeleted
		gitTagger := GitTagger{}
		filePath := "../../tests/terraform/resources/taggedkms/deleted_kms.tf"
		blame := taggedkms.KmsBlame
		gitTagger.mapOriginFileToGitFile(filePath, &blame)
		assert.Equal(t, expectedMapping["originToGit"], gitTagger.fileLinesMapper[filePath].originToGit)
		assert.Equal(t, expectedMapping["gitToOrigin"], gitTagger.fileLinesMapper[filePath].gitToOrigin)
	})
}
