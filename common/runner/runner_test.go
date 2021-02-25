package runner

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/tagging"
	terraformStructure "bridgecrewio/yor/terraform/structure"
	"bridgecrewio/yor/tests/utils/blameutils"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/stretchr/testify/assert"
)

func Test_loadExternalTags(t *testing.T) {
	t.Run("load local plugins", func(t *testing.T) {
		pluginDir := "../../tests/yor_plugins/example"
		fmt.Printf("please make sure you have .so file in %s. if not, run the following command: \n", pluginDir)
		fmt.Printf("go build -gcflags=\"all=-N -l\" -buildmode=plugin -o %s/extra_tags.so %s/*.go\n", pluginDir, pluginDir)
		gotTags, err := loadExternalTags([]string{pluginDir})
		if err != nil {
			t.Errorf("loadExternalTags() error = %v", err)
			return
		}
		expectedTags := map[string]string{"yor_checkov": "checkov", "git_owner": "bana"}
		assert.Equal(t, len(expectedTags), len(gotTags))
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		gitBlame := gitservice.GitBlame{
			GitOrg:        "bridgecrewio",
			GitRepository: "yor",
			BlamesByLine: map[int]*git.Line{0: {
				Author: "bana",
				Date:   now,
				Hash:   plumbing.NewHash("0"),
			}, 1: {Author: "shati",
				Date: yesterday,
				Hash: plumbing.NewHash("1")}}}
		for _, tag := range gotTags {
			tag.Init()
			tag, err := tag.CalculateValue(&gitBlame)
			print(err)
			key := tag.GetKey()
			value := tag.GetValue()
			assert.Equal(t, expectedTags[key], value)
		}
	})
}

func Test_E2E(t *testing.T) {
	t.Run("modified file not changing", func(t *testing.T) {
		filePath := "../../tests/terraform/resources/taggedkms/modified/modified_kms.tf"
		taggedFilePath := "../../tests/terraform/resources/taggedkms/modified/modified_kms_tagged.tf"

		defer func() {
			err := os.Remove(taggedFilePath)
			if err != nil {
				panic(err)
			}
		}()

		textBefore, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Errorf(fmt.Sprintf("Failed to read file %s because %s", filePath, err))
		}
		rootDir := "../../tests/terraform/resources/taggedkms/modified"
		gitTagger := initMockGitTagger(rootDir, map[string]string{filePath: "../../tests/terraform/resources/taggedkms/origin_kms.tf"})
		terraformParser := terraformStructure.TerrraformParser{}
		terraformParser.Init(rootDir, nil)

		blocks, err := terraformParser.ParseFile(filePath)
		if err != nil {
			t.Errorf(fmt.Sprintf("Failed to parse file %v", filePath))
		}
		for _, block := range blocks {
			if block.IsBlockTaggable() {
				gitTagger.CreateTagsForBlock(block)
			}
		}

		err = terraformParser.WriteFile(filePath, blocks, taggedFilePath)
		if err != nil {
			t.Errorf(fmt.Sprintf("Failed to write file %s because %s", taggedFilePath, err))
		}

		textAfter, err := ioutil.ReadFile(taggedFilePath)
		if err != nil {
			t.Errorf(fmt.Sprintf("Failed to read file %s because %s", taggedFilePath, err))
		}
		assert.Equal(t, textBefore, textAfter)
	})
}

func Test_TagCFNDir(t *testing.T) {
	t.Run("tag cloudformation yaml dir", func(t *testing.T) {
		options := common.Options{
			Directory: "../../tests/cloudformation/resources/ebs",
			ExtraTags: "{}",
		}
		filePath := options.Directory + "/ebs.yaml"

		originFileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		originFileLines := common.GetLinesFromBytes(originFileBytes)

		defer func() {
			_ = ioutil.WriteFile(filePath, originFileBytes, 0644)
		}()

		mockGitTagger := initMockGitTagger(options.Directory, map[string]string{filePath: filePath})
		runner := Runner{}
		err = runner.Init(&options)
		if err != nil {
			t.Error(err)
		}
		runner.taggers[0] = mockGitTagger
		_, err = runner.TagDirectory(options.Directory)
		if err != nil {
			t.Error(err)
		}

		editedFileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		editedFileLines := common.GetLinesFromBytes(editedFileBytes)

		expectedAddedLines := len(mockGitTagger.Tags) * 2
		assert.Equal(t, len(originFileLines)+expectedAddedLines, len(editedFileLines))

		matcher := difflib.NewMatcher(originFileLines, editedFileLines)
		matches := matcher.GetMatchingBlocks()
		expectedMatches := []difflib.Match{
			{A: 0, B: 0, Size: 13}, {A: 13, B: 29, Size: 2}, {A: 15, B: 31, Size: 0},
		}
		assert.Equal(t, expectedMatches, matches)
	})
}

func initMockGitTagger(rootDir string, filesToBlames map[string]string) *tagging.GitTagger {
	gitService, _ := gitservice.NewGitService(rootDir)

	for filePath := range filesToBlames {
		blameSrc, _ := ioutil.ReadFile(filesToBlames[filePath])
		blame := blameutils.CreateMockBlame(blameSrc)
		gitService.BlameByFile[filePath] = &blame
	}

	gitTagger := tagging.GitTagger{GitService: gitService}
	gitTagger.InitTags(nil)

	return &gitTagger
}
