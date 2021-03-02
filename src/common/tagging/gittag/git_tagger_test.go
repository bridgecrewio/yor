package gittag

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/gitservice"
	commonStructure "bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/tests/utils/blameutils"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestGitTagger(t *testing.T) {
	path := "../../../../tests/utils/blameutils/git_tagger_file.txt"
	blame := blameutils.SetupBlameResults(t, path, 3)

	t.Run("test git tagger CreateTagsForBlock", func(t *testing.T) {
		gitService := &gitservice.GitService{
			BlameByFile: map[string]*git.BlameResult{path: blame},
		}
		tagger := Tagger{}

		wd, _ := os.Getwd()
		tagger.InitTagger(wd, nil)
		tagger.GitService = gitService
		block := &MockTestBlock{
			Block: commonStructure.Block{
				FilePath:   path,
				IsTaggable: true,
			},
		}

		tagger.CreateTagsForBlock(block)
		assert.Equal(t, 7, len(block.NewTags))
	})
}

func TestGitTagger_mapOriginFileToGitFile(t *testing.T) {
	t.Run("map tagged kms", func(t *testing.T) {
		expectedMapping := ExpectedFileMappingTagged
		gitTagger := Tagger{}
		filePath := "../../../../tests/terraform/resources/taggedkms/tagged_kms.tf"
		src, _ := ioutil.ReadFile("../../../../tests/terraform/resources/taggedkms/origin_kms.tf")
		blame := blameutils.CreateMockBlame(src)
		gitTagger.mapOriginFileToGitFile(filePath, &blame)
		assert.Equal(t, expectedMapping["originToGit"], gitTagger.fileLinesMapper[filePath].originToGit)
		assert.Equal(t, expectedMapping["gitToOrigin"], gitTagger.fileLinesMapper[filePath].gitToOrigin)
	})
	t.Run("map kms with deleted lines", func(t *testing.T) {
		expectedMapping := ExpectedFileMappingDeleted
		gitTagger := Tagger{}
		filePath := "../../../../tests/terraform/resources/taggedkms/deleted_kms.tf"
		src, _ := ioutil.ReadFile("../../../../tests/terraform/resources/taggedkms/origin_kms.tf")
		blame := blameutils.CreateMockBlame(src)
		gitTagger.mapOriginFileToGitFile(filePath, &blame)
		assert.Equal(t, expectedMapping["originToGit"], gitTagger.fileLinesMapper[filePath].originToGit)
		assert.Equal(t, expectedMapping["gitToOrigin"], gitTagger.fileLinesMapper[filePath].gitToOrigin)
	})
}

var ExpectedFileMappingTagged = map[string]map[int]int{
	"originToGit": {1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: -1, 7: -1, 8: -1, 9: -1, 10: -1, 11: -1, 12: -1, 13: -1, 14: -1, 15: -1, 16: 6, 17: 7, 18: 8, 19: 9, 20: 10, 21: 11, 22: 12},
	"gitToOrigin": {1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 16, 7: 17, 8: 18, 9: 19, 10: 20, 11: 21, 12: 22},
}
var ExpectedFileMappingDeleted = map[string]map[int]int{
	"originToGit": {1: 1, 2: 2, 3: 3, 4: 4, 5: 6, 6: 7, 7: 8, 8: 9, 9: 10, 10: 11, 11: 12},
	"gitToOrigin": {1: 1, 2: 2, 3: 3, 4: 4, 5: -1, 6: 5, 7: 6, 8: 7, 9: 8, 10: 9, 11: 10, 12: 11},
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

func (b *MockTestBlock) GetLines(_ ...bool) common.Lines {
	return common.Lines{Start: 1, End: 3}
}
