package integration

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/runner"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMultipleCommits(t *testing.T) {
	t.Run("Test terragoat tagging", func(t *testing.T) {
		part1Text, err := ioutil.ReadFile("./resources/commits_file_1.tf")
		panicIfErr(err)

		part2Text, err := ioutil.ReadFile("./resources/commits_file_2.tf")
		panicIfErr(err)

		dir, err := ioutil.TempDir("", "temp-repo")
		panicIfErr(err)

		tfFileName := "main.tf"
		tfFilePath := path.Join(dir, tfFileName)
		err = ioutil.WriteFile(tfFilePath, part1Text, 0644)
		panicIfErr(err)

		testRepo, err := git.PlainInit(dir, false)
		panicIfErr(err)

		worktree, err := testRepo.Worktree()
		panicIfErr(err)

		_, err = worktree.Add(tfFileName)
		panicIfErr(err)

		commit1, err := worktree.Commit("commit resource 1 without tags", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana1",
				Email: "Bana1@gmail.com",
				When:  time.Now().AddDate(0, 0, -2),
			},
		})
		panicIfErr(err)

		yorRunner := runner.Runner{}
		err = yorRunner.Init(&common.Options{
			Directory: dir,
			ExtraTags: "{}",
		})
		panicIfErr(err)

		reportService2, err := yorRunner.TagDirectory(dir)
		panicIfErr(err)

		reportService2.CreateReport()
		report2 := reportService2.GetReport()
		var resource1Trace string
		for _, tag := range report2.NewResourceTags {
			if tag.TagKey == "git_commit" {
				assert.Equal(t, commit1.String(), tag.UpdatedValue)
			} else if tag.TagKey == "yor_trace" {
				resource1Trace = tag.UpdatedValue
			}
		}

		_, err = worktree.Add(tfFileName)
		panicIfErr(err)
		commit2, err := worktree.Commit("commit resource 1 added yor tags", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana2",
				Email: "Bana2@gmail.com",
				When:  time.Now().AddDate(0, 0, -1),
			},
		})

		f, err := os.OpenFile(tfFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		panicIfErr(err)

		defer f.Close()

		if _, err = f.Write(part2Text); err != nil {
			panic(err)
		}

		_, err = worktree.Add(tfFileName)
		panicIfErr(err)
		commit3, err := worktree.Commit("commit resource 2 without tags", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana3",
				Email: "Bana3@gmail.com",
				When:  time.Now(),
			},
		})
		yorRunner2 := runner.Runner{}
		err = yorRunner2.Init(&common.Options{
			Directory: dir,
			ExtraTags: "{}",
		})
		time.Sleep(2 * time.Second)
		reportService2, err = yorRunner2.TagDirectory(dir)
		panicIfErr(err)
		reportService2.CreateReport()
		fmt.Printf("Commit1: %s\nCommit2: %s\nCommit3: %s\n", commit1.String(), commit2.String(), commit3.String())
		report2 = reportService2.GetReport()
		reportService2.PrintToStdout()
		for _, tag := range report2.NewResourceTags {
			if tag.TagKey == "git_commit" && tag.ResourceID == "aws_s3_bucket.f2" {
				assert.Equal(t, commit3.String(), tag.UpdatedValue, "new resource should have commit3")
			}
		}
		for _, tag := range report2.UpdatedResourceTags {
			if tag.TagKey == "git_commit" && tag.ResourceID == "aws_s3_bucket.financials" {
				assert.Equal(t, commit2.String(), tag.UpdatedValue, "updated resource should be commit2")
			} else if tag.TagKey == "yor_trace" && tag.ResourceID == "aws_s3_bucket.financials" {
				assert.Equal(t, resource1Trace, tag.UpdatedValue)
			}
		}

	})
}
