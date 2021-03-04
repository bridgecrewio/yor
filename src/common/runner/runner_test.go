package runner

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/gittag"
	terraformStructure "bridgecrewio/yor/src/terraform/structure"
	"bridgecrewio/yor/tests/utils"
	"bridgecrewio/yor/tests/utils/blameutils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/stretchr/testify/assert"
)

func Test_loadExternalTags(t *testing.T) {
	t.Run("load tags plugins", func(t *testing.T) {
		pluginDir := "../../../tests/yor_plugins/example"
		fmt.Printf("please make sure you have .so file in %s. if not, run the following command: \n", pluginDir)
		fmt.Printf("go build -gcflags=\"all=-N -l\" -buildmode=plugin -o %s/extra_tags.so %s/*.go\n", pluginDir, pluginDir)
		gotTags, _, err := loadExternalResources([]string{pluginDir})
		if err != nil {
			t.Errorf("loadExternalResources() error = %v", err)
			return
		}
		expectedTags := map[string]string{"yor_foo": "foo", "git_owner": "bana"}
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
			tagVal, err := tag.CalculateValue(&gitBlame)
			print(err)
			key := tagVal.GetKey()
			value := tagVal.GetValue()
			assert.Equal(t, expectedTags[key], value)
		}
	})

	t.Run("load taggers plugins", func(t *testing.T) {
		pluginDir := "../../../tests/yor_plugins/tagger_example"
		fmt.Printf("please make sure you have .so file in %s. if not, run the following command: \n", pluginDir)
		fmt.Printf("go build -gcflags=\"all=-N -l\" -buildmode=plugin -o %s/extra_tags.so %s/*.go\n", pluginDir, pluginDir)
		gotTags, _, err := loadExternalResources([]string{pluginDir})
		if err != nil {
			t.Errorf("loadExternalResources() error = %v", err)
			return
		}
		expectedTags := map[string]string{"bc_dir": "tests/yor_plugins/tagger_example"}
		assert.Equal(t, len(expectedTags), len(gotTags))
		for _, tag := range gotTags {
			tag.Init()
			tagVal, _ := tag.CalculateValue(&structure.Block{FilePath: "some/path/to/file.tf"})
			key := tagVal.GetKey()
			value := tagVal.GetValue()
			assert.Equal(t, expectedTags[key], value)
		}
	})
}

func Test_E2E(t *testing.T) {
	t.Run("modified file not changing", func(t *testing.T) {
		filePath := "../../../tests/terraform/resources/taggedkms/modified/modified_kms.tf"
		taggedFilePath := "../../../tests/terraform/resources/taggedkms/modified/modified_kms_tagged.tf"

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
		rootDir := "../../../tests/terraform/resources/taggedkms/modified"
		gitTagger := initMockGitTagger(rootDir, map[string]string{filePath: "../../../tests/terraform/resources/taggedkms/origin_kms.tf"})
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
	t.Run("tag cloudformation yaml with tags", func(t *testing.T) {
		options := common.Options{
			Directory: "../../../tests/cloudformation/resources/ebs",
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
		_, err = runner.TagDirectory()
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)

		editedFileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		editedFileLines := common.GetLinesFromBytes(editedFileBytes)

		expectedAddedLines := len(mockGitTagger.GetTags())*2 + 2
		assert.Equal(t, len(originFileLines)+expectedAddedLines, len(editedFileLines))

		matcher := difflib.NewMatcher(originFileLines, editedFileLines)
		matches := matcher.GetMatchingBlocks()
		expectedMatches := []difflib.Match{
			{A: 0, B: 0, Size: 13}, {A: 13, B: 29, Size: 2}, {A: 15, B: 31, Size: 0},
		}
		assert.Equal(t, expectedMatches, matches)
	})
}

func TestRunnerInternals(t *testing.T) {
	t.Run("Test isFileSkipped", func(t *testing.T) {
		runner := Runner{}
		rootDir := "../../../tests/terraform"
		skippedFiles := []string{"../../../tests/terraform/mixed/mixed.tf", "../../../tests/terraform/resources/tagged/complex_tags_tagged.tf"}
		_ = runner.Init(&common.Options{
			Directory: rootDir,
			SkipDirs:  []string{"../../../tests/terraform/mixed", "../../../tests/terraform/resources/tagged/"},
			ExtraTags: "{}",
		})

		_ = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				isFileSkipped := runner.isFileSkipped(&terraformStructure.TerrraformParser{}, path)
				if isFileSkipped {
					shouldSkip := false
					skippedIndex := -1
					for i, skipped := range skippedFiles {
						if skipped == path {
							shouldSkip = true
							skippedIndex = i
						}
					}
					if shouldSkip {
						skippedFiles = append(skippedFiles[:skippedIndex], skippedFiles[skippedIndex+1:]...)
					} else if strings.HasSuffix(path, ".tf") {
						assert.Fail(t, fmt.Sprintf("Should not have skipped %v", path))
					}
				}
			}
			return nil
		})

		assert.Equal(t, 0, len(skippedFiles), "Some files were not skipped")
	})

	t.Run("Test skip entire dir", func(t *testing.T) {
		runner := Runner{}
		rootDir := "../../../tests/terraform"
		output := utils.CaptureOutput(func() {
			_ = runner.Init(&common.Options{
				Directory: rootDir,
				SkipDirs: []string{
					"../../../tests/terraform/mixed",
					"../../../tests/terraform/resources/tagged/",
					"../../../tests/terraform",
				},
				ExtraTags: "{}",
			})
		})
		assert.Contains(t, output, "[WARNING] Selected dir, ../../../tests/terraform, is skipped - expect an empty result")
	})
}

func initMockGitTagger(rootDir string, filesToBlames map[string]string) *gittag.Tagger {
	gitService, _ := gitservice.NewGitService(rootDir)

	for filePath := range filesToBlames {
		blameSrc, _ := ioutil.ReadFile(filesToBlames[filePath])
		blame := blameutils.CreateMockBlame(blameSrc)
		gitService.BlameByFile[filePath] = &blame
	}

	gitTagger := gittag.Tagger{}
	wd, _ := os.Getwd()
	gitTagger.InitTagger(wd, nil)
	gitTagger.GitService = gitService
	return &gitTagger
}
