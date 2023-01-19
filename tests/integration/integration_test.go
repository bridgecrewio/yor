package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/bridgecrewio/yor/src/common/clioptions"
	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/reports"
	"github.com/bridgecrewio/yor/src/common/runner"
	"github.com/bridgecrewio/yor/src/common/tagging/gittag"
	tagUtils "github.com/bridgecrewio/yor/src/common/tagging/utils"
	terraformStructure "github.com/bridgecrewio/yor/src/terraform/structure"
	"github.com/bridgecrewio/yor/tests/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
)

func TestMultipleCommits(t *testing.T) {
	t.Run("Test tagging over multiple commits", func(t *testing.T) {
		// read two resource files to be added to a new file we create
		part1Text, err := ioutil.ReadFile("./resources/commits_file_1.tf")
		failIfErr(t, err)
		part2Text, err := ioutil.ReadFile("./resources/commits_file_2.tf")
		failIfErr(t, err)

		// init temp directory and file, and write the first text to it
		dir, err := ioutil.TempDir("", "commits")
		failIfErr(t, err)
		defer func() {
			_ = os.RemoveAll(dir)
		}()
		tfFileName := "main.tf"
		tfFilePath := path.Join(dir, tfFileName)
		err = ioutil.WriteFile(tfFilePath, part1Text, 0644)
		failIfErr(t, err)

		// init git repository and commit the file
		testRepo, err := git.PlainInit(dir, false)
		failIfErr(t, err)
		worktree, err := testRepo.Worktree()
		failIfErr(t, err)
		commit1 := commitFile(worktree, tfFileName, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana1",
				Email: "Bana1@gmail.com",
				When:  time.Now().AddDate(0, 0, -2),
			},
		})
		time.Sleep(2 * time.Second)
		// run yor on resource 1
		yorRunner := runner.Runner{}
		err = yorRunner.Init(&clioptions.TagOptions{
			Directory: dir,
			TagGroups: getTagGroups(),
		})
		failIfErr(t, err)
		reportService, err := yorRunner.TagDirectory()
		failIfErr(t, err)
		reportService.CreateReport()
		report := reportService.GetReport()

		// check if the resource has the right commit hash and save the yor trace
		var resource1Trace string
		for _, tag := range report.NewResourceTags {
			if tag.TagKey == "git_commit" {
				assert.Equal(t, commit1.String(), tag.UpdatedValue)
			} else if tag.TagKey == "yor_trace" {
				resource1Trace = tag.UpdatedValue
			}
		}

		// commit the added tags
		commit2 := commitFile(worktree, tfFileName, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana2",
				Email: "Bana2@gmail.com",
				When:  time.Now().AddDate(0, 0, -1),
			},
		})
		time.Sleep(2 * time.Second)

		// append to the file the second resource
		f, err := os.OpenFile(tfFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		failIfErr(t, err)
		defer func() {
			_ = f.Close()
		}()
		if _, err = f.Write(part2Text); err != nil {
			panic(err)
		}

		// commit the second resource
		commit3 := commitFile(worktree, tfFileName, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana3",
				Email: "Bana3@gmail.com",
				When:  time.Now(),
			},
		})

		// run yor on both resources
		yorRunner2 := runner.Runner{}
		err = yorRunner2.Init(&clioptions.TagOptions{
			Directory: dir,
			TagGroups: getTagGroups(),
		})
		failIfErr(t, err)
		time.Sleep(2 * time.Second)
		reportService, err = yorRunner2.TagDirectory()
		failIfErr(t, err)
		reportService.CreateReport()
		report2 := reportService.GetReport()

		// check if the second resource has the third commit
		for _, tag := range report2.NewResourceTags {
			if tag.TagKey == "git_commit" && tag.ResourceID == "aws_s3_bucket.f2" {
				assert.Equal(t, commit3.String(), tag.UpdatedValue, "new resource should have commit3")
			}
		}

		// check if the first resource has the second commit (because we committed the tags) and that the trace hasn't changed
		for _, tag := range report2.UpdatedResourceTags {
			if tag.TagKey == "git_commit" && tag.ResourceID == "aws_s3_bucket.financials" {
				assert.Equal(t, commit2.String(), tag.UpdatedValue, "updated resource should be commit2")
			} else if tag.TagKey == "yor_trace" && tag.ResourceID == "aws_s3_bucket.financials" {
				assert.Equal(t, resource1Trace, tag.UpdatedValue)
			}
		}

	})
}

func TestRunResults(t *testing.T) {
	t.Run("Test terragoat tagging", func(t *testing.T) {
		content, _ := ioutil.ReadFile("../../result.json")
		report := &reports.Report{}
		err := json.Unmarshal(content, &report)
		if err != nil {
			assert.Fail(t, "Failed to parse json result")
		}
		assert.Less(t, 63, report.Summary.Scanned)
		assert.LessOrEqual(t, 63, report.Summary.NewResources)
		assert.Equal(t, 0, report.Summary.UpdatedResources)

		var taggedAWS, taggedGCP, taggedAzure int
		resourceSet := make(map[string]bool)

		for _, tr := range report.NewResourceTags {
			switch {
			case strings.HasPrefix(tr.ResourceID, "aws"):
				taggedAWS++
			case strings.HasPrefix(tr.ResourceID, "google_"):
				taggedGCP++
			case strings.HasPrefix(tr.ResourceID, "azurerm"):
				taggedAzure++
			}

			assert.NotEqual(t, "", tr.ResourceID)
			assert.NotEqual(t, "", tr.File)
			assert.NotEqual(t, "", tr.UpdatedValue)
			assert.NotEqual(t, "", tr.TagKey)
			assert.NotEqual(t, "", tr.YorTraceID)
			assert.Equal(t, "", tr.OldValue)

			resourceSet[tr.ResourceID] = true
		}

		assert.LessOrEqual(t, 312, taggedAWS)
		assert.LessOrEqual(t, 32, taggedGCP)
		assert.LessOrEqual(t, 160, taggedAzure)
		assert.Equal(t, report.Summary.NewResources, len(resourceSet))
	})

	t.Run("Test cli arg parsing", func(t *testing.T) {
		resultFile := "../../list-tags-result.txt"
		content, _ := ioutil.ReadFile(resultFile)
		defer func() {
			_ = os.Remove(resultFile)
		}()
		lines := strings.Split(string(content), "\n")
		var filtered []string
		for _, line := range lines {
			if len(line) > 3 && line[2] != ' ' && line[2] != '-' {
				filtered = append(filtered, line)
			}
		}
		if len(filtered) == 2 {
			assert.True(t, (strings.HasPrefix(filtered[0], "| code2cloud") && strings.HasPrefix(filtered[1], "| git")) ||
				(strings.HasPrefix(filtered[0], "| git") && strings.HasPrefix(filtered[1], "| code2cloud")))
		} else {
			assert.Fail(t, fmt.Sprintf("Number of filtered lines is %v, should be %v", len(filtered), 2))
		}
	})

	t.Run("Test terraform-aws-bridgecrew-read-only tagging specified tags", func(t *testing.T) {
		repoPath := utils.CloneRepo("https://github.com/bridgecrewio/terraform-aws-bridgecrew-read-only.git", "a8686215642fd47a38bf8615d91d0d40630ab989")
		defer os.RemoveAll(repoPath)

		yorRunner := new(runner.Runner)
		err := yorRunner.Init(&clioptions.TagOptions{
			Directory: repoPath,
			TagGroups: getTagGroups(),
			Tag:       []string{"yor_trace"},
			Parsers:   []string{"Terraform"},
		})
		failIfErr(t, err)
		reportService, err := yorRunner.TagDirectory()
		failIfErr(t, err)

		reportService.CreateReport()
		report := reportService.GetReport()
		assert.LessOrEqual(t, 18, report.Summary.Scanned)
		assert.Greater(t, report.Summary.Scanned, 0)

		for _, newTag := range report.NewResourceTags {
			if strings.HasPrefix(repoPath, newTag.File) {
				assert.Equal(t, "yor_trace", newTag.TagKey)
				assert.Equal(t, "aws_iam_role.bridgecrew_account_role", newTag.ResourceID)
			}
		}
	})

	t.Run("Test terraform-aws-bridgecrew-read-only tagging skip tags", func(t *testing.T) {
		repoPath := utils.CloneRepo("https://github.com/bridgecrewio/terraform-aws-bridgecrew-read-only.git", "a8686215642fd47a38bf8615d91d0d40630ab989")
		defer os.RemoveAll(repoPath)

		yorRunner := runner.Runner{}
		err := yorRunner.Init(&clioptions.TagOptions{
			Directory: repoPath,
			TagGroups: getTagGroups(),
			SkipTags:  []string{"yor_trace"},
			Parsers:   []string{"Terraform"},
		})
		failIfErr(t, err)
		reportService, err := yorRunner.TagDirectory()
		failIfErr(t, err)

		reportService.CreateReport()
		report := reportService.GetReport()

		newTags := report.NewResourceTags
		for _, newTag := range newTags {
			if strings.HasPrefix(repoPath, newTag.File) {
				assert.NotEqual(t, "yor_trace", newTag.TagKey)
			}
		}
	})
}

func TestTagUncommittedResults(t *testing.T) {
	t.Run("Test tagging twice no result second time", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		outputPath := "./result_uncommitted.json"
		defer func() {
			_ = os.RemoveAll(terragoatPath)
			_ = os.RemoveAll(outputPath)
		}()

		terragoatAWSDirectory := path.Join(terragoatPath, "terraform/aws")

		// tag aws directory
		tagDirectory(t, terragoatAWSDirectory)
		// tag again, this time the files have uncommitted changes
		tagDirectory(t, terragoatAWSDirectory)

		terraformParser := terraformStructure.TerraformParser{}
		terraformParser.Init(terragoatAWSDirectory, nil)

		dbAppFile := path.Join(terragoatAWSDirectory, "db-app.tf")
		blocks, err := terraformParser.ParseFile(dbAppFile)
		failIfErr(t, err)
		defaultInstanceBlock := blocks[0].(*terraformStructure.TerraformBlock)
		if defaultInstanceBlock.GetResourceID() != "aws_db_instance.default" {
			t.Errorf("invalid file structure, the resource id is %s", defaultInstanceBlock.GetResourceID())
		}

		rawTags := defaultInstanceBlock.HclSyntaxBlock.Body.Attributes["tags"]
		rawTagsExpr := rawTags.Expr.(*hclsyntax.ObjectConsExpr)
		assert.Equal(t, "tags", rawTags.Name)
		assert.Equal(t, 10, len(rawTagsExpr.Items))

		currentTags := defaultInstanceBlock.ExitingTags

		expectedTagsValues := map[string]string{
			"Name":                 "${local.resource_prefix.value}-rds",
			"Environment":          "local.resource_prefix.value",
			"git_last_modified_by": "nimrodkor@gmail.com",
			"git_commit":           "d68d2897add9bc2203a5ed0632a5cdd8ff8cefb0",
			"git_file":             strings.TrimPrefix(dbAppFile, terragoatPath+"/"),
		}

		for _, tag := range currentTags {
			if expectedVal, ok := expectedTagsValues[tag.GetKey()]; ok {
				assert.Equal(t, expectedVal, tag.GetValue(), fmt.Sprintf("Missmach in tag %s, expected %s, got %s", tag.GetKey(), expectedVal, tag.GetValue()))
			}
			if tag.GetKey() == "git_last_modified_at" {
				timeTagValue, err := time.Parse("2006-01-02 15:04:05", tag.GetValue())
				failIfErr(t, err)
				diff := time.Now().UTC().Sub(timeTagValue)
				assert.Greater(t, diff.Hours(), 100.)
			}
		}
	})

	t.Run("Test tagging after minor change", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		outputPath := "./result_uncommitted.json"
		defer func() {
			_ = os.RemoveAll(terragoatPath)
			_ = os.RemoveAll(outputPath)
		}()

		terragoatAWSDirectory := path.Join(terragoatPath, "terraform/aws")

		// tag aws directory
		tagDirectory(t, terragoatAWSDirectory)

		// Make minor change to file
		input, _ := ioutil.ReadFile(path.Join(terragoatAWSDirectory, "db-app.tf"))
		lines := strings.Split(string(input), "\n")
		for i, line := range lines {
			if line == "  instance_class          = \"db.t3.micro\"" {
				lines[i] = "  instance_class          = \"db.t3.medium\""
			}
		}
		output := strings.Join(lines, "\n")
		_ = ioutil.WriteFile(path.Join(terragoatAWSDirectory, "db-app.tf"), []byte(output), 0644)

		// tag again, this time the files have uncommitted changes
		tagDirectory(t, terragoatAWSDirectory)

		terraformParser := terraformStructure.TerraformParser{}
		terraformParser.Init(terragoatAWSDirectory, nil)

		dbAppFile := path.Join(terragoatAWSDirectory, "db-app.tf")
		blocks, err := terraformParser.ParseFile(dbAppFile)
		failIfErr(t, err)
		defaultInstanceBlock := blocks[0].(*terraformStructure.TerraformBlock)
		if defaultInstanceBlock.GetResourceID() != "aws_db_instance.default" {
			t.Errorf("invalid file structure, the resource id is %s", defaultInstanceBlock.GetResourceID())
		}

		rawTags := defaultInstanceBlock.HclSyntaxBlock.Body.Attributes["tags"]
		rawTagsExpr := rawTags.Expr.(*hclsyntax.ObjectConsExpr)
		assert.Equal(t, "tags", rawTags.Name)
		assert.Equal(t, 10, len(rawTagsExpr.Items))

		currentTags := defaultInstanceBlock.ExitingTags

		expectedTagsValues := map[string]string{
			"Name":                 "${local.resource_prefix.value}-rds",
			"Environment":          "local.resource_prefix.value",
			"git_last_modified_by": gitservice.GetGitUserEmail(),
			"git_commit":           gittag.CommitUnavailable,
			"git_file":             strings.TrimPrefix(dbAppFile, terragoatPath+"/"),
		}

		for _, tag := range currentTags {
			if expectedVal, ok := expectedTagsValues[tag.GetKey()]; ok {
				assert.Equal(t, expectedVal, tag.GetValue(), fmt.Sprintf("Missmach in tag %s, expected %s, got %s", tag.GetKey(), expectedVal, tag.GetValue()))
			}
			if tag.GetKey() == "git_last_modified_at" {
				timeTagValue, err := time.Parse("2006-01-02 15:04:05", tag.GetValue())
				failIfErr(t, err)
				diff := time.Now().UTC().Sub(timeTagValue)
				assert.True(t, diff < 2*time.Minute)
			}
		}
	})
}

func TestLocalModules(t *testing.T) {
	t.Run("Test tagging local modules", func(t *testing.T) {
		localTagExampleRepo := "https://github.com/JamesWoolfenden/terraform-aws-activemq/"
		repoPath := utils.CloneRepo(localTagExampleRepo, "05ab598c4947bb9e53fee67a0f7350941897c2bd")
		defer func() {
			_ = os.RemoveAll(repoPath)
		}()

		targetDirectory := path.Join(repoPath, "example/examplea")

		tagLocalDirectory(t, targetDirectory)

		terraformParser := terraformStructure.TerraformParser{}
		terraformParser.Init(targetDirectory, nil)
		dbAppFile := path.Join(targetDirectory, "module.broker.tf")
		blocks, _ := terraformParser.ParseFile(dbAppFile)

		defaultInstanceBlock := blocks[0].(*terraformStructure.TerraformBlock)
		rawTags := defaultInstanceBlock.HclSyntaxBlock.Body.Attributes["common_tags"]
		rawTagsExpr := rawTags.Expr.(*hclsyntax.FunctionCallExpr)
		assert.Equal(t, "common_tags", rawTags.Name)
		assert.Equal(t, "merge", rawTagsExpr.Name)
	})

}

func failIfErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func tagDirectory(t *testing.T, path string) {
	yorRunner := runner.Runner{}
	err := yorRunner.Init(&clioptions.TagOptions{
		Directory: path,
		TagGroups: getTagGroups(),
		Parsers:   []string{"Terraform", "CloudFormation", "Serverless"},
	})
	failIfErr(t, err)
	_, err = yorRunner.TagDirectory()
	failIfErr(t, err)
}

func tagLocalDirectory(t *testing.T, path string) {
	yorRunner := runner.Runner{}
	err := yorRunner.Init(&clioptions.TagOptions{
		Directory:       path,
		TagGroups:       getTagGroups(),
		Parsers:         []string{"Terraform"},
		TagLocalModules: true,
	})
	failIfErr(t, err)
	_, err = yorRunner.TagDirectory()
	failIfErr(t, err)
}

func commitFile(worktree *git.Worktree, filename string, commitOptions *git.CommitOptions) plumbing.Hash {
	_, err := worktree.Add(filename)
	if err != nil {
		panic(err)
	}
	commit, err := worktree.Commit("commit resource 1 without tags", commitOptions)
	if err != nil {
		panic(err)
	}
	return commit
}

func getTagGroups() []string {
	return tagUtils.GetAllTagGroupsNames()
}
