package runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cloudformationStructure "github.com/bridgecrewio/yor/src/cloudformation/structure"
	"github.com/bridgecrewio/yor/src/common/cli"
	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/gittag"
	taggingUtils "github.com/bridgecrewio/yor/src/common/tagging/utils"
	"github.com/bridgecrewio/yor/src/common/utils"
	terraformStructure "github.com/bridgecrewio/yor/src/terraform/structure"
	testingUtils "github.com/bridgecrewio/yor/tests/utils"
	"github.com/bridgecrewio/yor/tests/utils/blameutils"

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

	t.Run("load tagGroups plugins", func(t *testing.T) {
		pluginDir := "../../../tests/yor_plugins/tag_group_example"
		fmt.Printf("please make sure you have .so file in %s. if not, run the following command: \n", pluginDir)
		fmt.Printf("go build -gcflags=\"all=-N -l\" -buildmode=plugin -o %s/extra_tag_groups.so %s/*.go\n", pluginDir, pluginDir)
		_, gotTagGroups, err := loadExternalResources([]string{pluginDir})
		if err != nil {
			t.Errorf("loadExternalResources() error = %v", err)
			return
		}
		assert.Equal(t, 1, len(gotTagGroups))
		group := gotTagGroups[0]
		group.InitTagGroup("src", nil)
		groupTags := gotTagGroups[0].GetTags()
		assert.Equal(t, 1, len(gotTagGroups[0].GetTags()))
		tag := groupTags[0]
		assert.Equal(t, "custom_owner", tag.GetKey())
		tagVal, _ := tag.CalculateValue(&terraformStructure.TerraformBlock{Block: structure.Block{FilePath: "src/auth/index.js"}})
		assert.Equal(t, "custom_owner", tagVal.GetKey())
		assert.Equal(t, "team-infra@company.com", tagVal.GetValue())

		tagVal, _ = tag.CalculateValue(&cloudformationStructure.CloudformationBlock{Block: structure.Block{FilePath: "src/some/path"}})
		assert.Equal(t, "custom_owner", tagVal.GetKey())
		assert.Equal(t, "team-it@company.com", tagVal.GetValue())
	})
}

func Test_TagCFNDir(t *testing.T) {
	t.Run("tag cloudformation yaml with tags", func(t *testing.T) {
		options := cli.TagOptions{
			Directory: "../../../tests/cloudformation/resources/ebs",
			TagGroups: taggingUtils.GetAllTagGroupsNames(),
		}
		filePath := options.Directory + "/ebs.yaml"

		originFileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		originFileLines := utils.GetLinesFromBytes(originFileBytes)

		defer func() {
			_ = ioutil.WriteFile(filePath, originFileBytes, 0644)
		}()

		mockGitTagGroup := initMockGitTagGroup(options.Directory, map[string]string{filePath: filePath})
		runner := Runner{}
		err = runner.Init(&options)
		if err != nil {
			t.Error(err)
		}
		runner.TagGroups[0] = mockGitTagGroup
		_, err = runner.TagDirectory()
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)

		editedFileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		editedFileLines := utils.GetLinesFromBytes(editedFileBytes)

		expectedAddedLines := len(mockGitTagGroup.GetTags()) * 2
		assert.Equal(t, len(originFileLines)+expectedAddedLines-1, len(editedFileLines))

		matcher := difflib.NewMatcher(originFileLines, editedFileLines)
		matches := matcher.GetMatchingBlocks()
		expectedMatches := []difflib.Match{
			{A: 0, B: 0, Size: 8}, {A: 8, B: 10, Size: 2}, {A: 11, B: 13, Size: 1}, {A: 15, B: 28, Size: 0},
		}
		assert.Equal(t, expectedMatches, matches)
	})

	t.Run("Filter tag groups", func(t *testing.T) {
		runner := Runner{}
		allTagGroups := taggingUtils.GetAllTagGroupsNames()
		_ = runner.Init(&cli.TagOptions{
			Directory: "../../../tests/cloudformation/resources/ebs",
			TagGroups: allTagGroups[:len(allTagGroups)-1],
		})

		tg := runner.TagGroups
		assert.Equal(t, len(allTagGroups)-1, len(tg))
	})
}

func TestRunnerInternals(t *testing.T) {
	t.Run("Test isFileSkipped", func(t *testing.T) {
		runner := Runner{}
		rootDir := "../../../tests/terraform"
		skippedFiles := []string{"../../../tests/terraform/mixed/mixed.tf", "../../../tests/terraform/resources/tagged/complex_tags_tagged.tf"}
		_ = runner.Init(&cli.TagOptions{
			Directory: rootDir,
			SkipDirs:  []string{"../../../tests/terraform/mixed", "../../../tests/terraform/resources/tagged/"},
			TagGroups: taggingUtils.GetAllTagGroupsNames(),
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
		output := testingUtils.CaptureOutput(func() {
			_ = runner.Init(&cli.TagOptions{
				Directory: rootDir,
				SkipDirs: []string{
					"../../../tests/terraform/mixed",
					"../../../tests/terraform/resources/tagged/",
					"../../../tests/terraform",
				},
				TagGroups: taggingUtils.GetAllTagGroupsNames(),
			})
		})
		assert.Contains(t, output, "[WARNING] Selected dir, ../../../tests/terraform, is skipped - expect an empty result")
	})
}

func initMockGitTagGroup(rootDir string, filesToBlames map[string]string) *gittag.TagGroup {
	gitService, _ := gitservice.NewGitService(rootDir)

	for filePath := range filesToBlames {
		blameSrc, _ := ioutil.ReadFile(filesToBlames[filePath])
		blame := blameutils.CreateMockBlame(blameSrc)
		gitService.BlameByFile[filePath] = &blame
	}

	gitTagGroup := gittag.TagGroup{}
	wd, _ := os.Getwd()
	gitTagGroup.InitTagGroup(wd, nil)
	gitTagGroup.GitService = gitService
	return &gitTagGroup
}
