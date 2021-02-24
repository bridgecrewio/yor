package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	commonStructure "bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"bridgecrewio/yor/tests/utils/blameutils"
	"bridgecrewio/yor/tests/utils/structureutils"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"

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
		expectedMapping := ExpectedFileMappingTagged
		gitTagger := GitTagger{}
		filePath := "../../tests/terraform/resources/taggedkms/tagged_kms.tf"
		blame := KmsBlame
		gitTagger.mapOriginFileToGitFile(filePath, &blame)
		assert.Equal(t, expectedMapping["originToGit"], gitTagger.fileLinesMapper[filePath].originToGit)
		assert.Equal(t, expectedMapping["gitToOrigin"], gitTagger.fileLinesMapper[filePath].gitToOrigin)
	})
	t.Run("map kms with deleted lines", func(t *testing.T) {
		expectedMapping := ExpectedFileMappingDeleted
		gitTagger := GitTagger{}
		filePath := "../../tests/terraform/resources/taggedkms/deleted_kms.tf"
		blame := KmsBlame
		gitTagger.mapOriginFileToGitFile(filePath, &blame)
		assert.Equal(t, expectedMapping["originToGit"], gitTagger.fileLinesMapper[filePath].originToGit)
		assert.Equal(t, expectedMapping["gitToOrigin"], gitTagger.fileLinesMapper[filePath].gitToOrigin)
	})
}

var layout = "2006-01-02 15:04:05"

func getTime() time.Time {
	t, _ := time.Parse(layout, "2020-06-16 17:46:24")
	return t
}

var KmsBlame = git.BlameResult{Lines: []*git.Line{
	{
		Author: "nimrodkor@gmail.com",
		Text:   "resource \"aws_kms_key\" \"logs_key\" {",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "  # key does not have rotation enabled",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "  description = \"${local.resource_prefix.value}-logs bucket key\"",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "  deletion_window_in_days = 7",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "}",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "resource \"aws_kms_alias\" \"logs_key_alias\" {",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "  name          = \"alias/${local.resource_prefix.value}-logs-bucket-key\"",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "  target_key_id = \"${aws_kms_key.logs_key.key_id}\"",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
	{
		Author: "nimrodkor@gmail.com",
		Text:   "}",
		Date:   getTime(),
		Hash:   plumbing.NewHash("d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0"),
	},
}}

var ExpectedFileMappingTagged = map[string]map[int]int{
	"originToGit": {1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: -1, 7: -1, 8: -1, 9: -1, 10: -1, 11: -1, 12: -1, 13: -1, 14: -1, 15: -1, 16: 6, 17: 7, 18: 8, 19: 9, 20: 10, 21: 11, 22: -1},
	"gitToOrigin": {1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 16, 7: 17, 8: 18, 9: 19, 10: 20, 11: 21},
}
var ExpectedFileMappingDeleted = map[string]map[int]int{
	"originToGit": {1: 1, 2: 2, 3: 3, 4: 4, 5: 6, 6: 7, 7: 8, 8: 9, 9: 10, 10: 11, 11: -1},
	"gitToOrigin": {1: 1, 2: 2, 3: 3, 4: 4, 5: -1, 6: 5, 7: 6, 8: 7, 9: 8, 10: 9, 11: 10},
}
