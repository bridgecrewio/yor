package integration

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/runner"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func commitFile(worktree *git.Worktree, filename string, commitOptions *git.CommitOptions) plumbing.Hash {
	_, err := worktree.Add(filename)
	panicIfErr(err)
	commit, err := worktree.Commit("commit resource 1 without tags", commitOptions)
	panicIfErr(err)
	return commit
}

func TestMultipleCommits(t *testing.T) {
	t.Run("Test terragoat tagging", func(t *testing.T) {
		// read two resource files to be added to a new file we create
		part1Text, err := ioutil.ReadFile("./resources/commits_file_1.tf")
		panicIfErr(err)
		part2Text, err := ioutil.ReadFile("./resources/commits_file_2.tf")
		panicIfErr(err)

		// init temp directory and file, and write the first text to it
		dir, err := ioutil.TempDir("", "temp-repo")
		panicIfErr(err)
		defer os.RemoveAll(dir)
		tfFileName := "main.tf"
		tfFilePath := path.Join(dir, tfFileName)
		err = ioutil.WriteFile(tfFilePath, part1Text, 0644)
		panicIfErr(err)

		// init git repository and commit the file
		testRepo, err := git.PlainInit(dir, false)
		panicIfErr(err)
		worktree, err := testRepo.Worktree()
		panicIfErr(err)
		commit1 := commitFile(worktree, tfFileName, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana1",
				Email: "Bana1@gmail.com",
				When:  time.Now().AddDate(0, 0, -2),
			},
		})

		// run yor on resource 1
		yorRunner := runner.Runner{}
		err = yorRunner.Init(&common.Options{
			Directory: dir,
			ExtraTags: "{}",
		})
		panicIfErr(err)
		reportService, err := yorRunner.TagDirectory()
		panicIfErr(err)
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
				When:  time.Now().AddDate(0, 0, -2),
			},
		})

		// append to the file the second resource
		f, err := os.OpenFile(tfFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		panicIfErr(err)
		defer f.Close()
		if _, err = f.Write(part2Text); err != nil {
			panic(err)
		}

		// commit the second resource
		commit3 := commitFile(worktree, tfFileName, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Bana3",
				Email: "Bana3@gmail.com",
				When:  time.Now().AddDate(0, 0, -2),
			},
		})

		// run yor on both resources
		yorRunner2 := runner.Runner{}
		err = yorRunner2.Init(&common.Options{
			Directory: dir,
			ExtraTags: "{}",
		})
		panicIfErr(err)
		time.Sleep(2 * time.Second)
		reportService, err = yorRunner2.TagDirectory()
		panicIfErr(err)
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
