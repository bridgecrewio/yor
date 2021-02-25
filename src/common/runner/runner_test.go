package runner

import (
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/tagging"
	terraformStructure "bridgecrewio/yor/src/terraform/structure"
	"fmt"
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
		pluginDir := "../../../tests/yor_plugins/example"
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
		blame := KmsBlame
		gitService, err := gitservice.NewGitService(rootDir)
		if err != nil {
			t.Errorf(fmt.Sprintf("Failed to init git service: %s", err))
		}
		gitService.BlameByFile = map[string]*git.BlameResult{filePath: &blame}
		gitTagger := tagging.GitTagger{GitService: gitService}
		gitTagger.InitTags(nil)
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
