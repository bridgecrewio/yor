package utils

import (
	"bridgecrewio/yor/common/git_service"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const Org = "bridgecrewio"
const Repository = "terragoat"
const FilePath = "README.md"
const CommitHash1 = "47accf06f13b503f3bab06fed7860e72f7523cac"
const CommitHash2 = "b2dc884b7439882c4dbe1e660cb1e02a3f84e45d"

func SetupBlame(t *testing.T) git_service.GitBlame {
	dateStr0 := "2020-03-28T21:42:46.000Z"
	dateStr1 := "2020-03-27T11:56:33.000Z"
	firstCommitDate, err1 := ExtractDate(dateStr0)
	secondCommitDate, err2 := ExtractDate(dateStr1)
	if err1 != nil || err2 != nil {
		assert.Fail(t, "Failed to parse static date")
	}
	return git_service.GitBlame{
		GitOrg:        Org,
		GitRepository: Repository,
		FilePath:      FilePath,
		BlamesByLine: map[int]*git.Line{
			0: {
				Author: "schosterbarak@gmail.com",
				Text:   "# Terragoat",
				Date:   firstCommitDate,
				Hash:   plumbing.NewHash(CommitHash1),
			},
			1: {
				Author: "jonjozwiak@users.noreply.github.com",
				Text:   "Bridgecrew solution to create vulnerable infrastructure",
				Date:   secondCommitDate,
				Hash:   plumbing.NewHash(CommitHash2),
			},
		},
	}
}

func ExtractDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	parsedDate, err := time.Parse(layout, dateStr)
	return parsedDate, err
}
