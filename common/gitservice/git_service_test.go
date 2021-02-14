package gitservice

import (
	"bridgecrewio/yor/tests"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TerragoatURL = "https://github.com/bridgecrewio/terragoat.git"

func extractUnixDateFromLine(line string) time.Time {
	userDateNum := line[strings.Index(line, "(")+1 : strings.Index(line, ")")]
	splitUserDateNum := strings.Split(userDateNum, " ")
	unixString := splitUserDateNum[len(splitUserDateNum)-2]
	intUnix, _ := strconv.ParseInt(unixString, 10, 64)

	return time.Unix(intUnix, 0)
}

func TestNewGitService(t *testing.T) {
	t.Run("Get correct organization and repo name", func(t *testing.T) {
		terragoatPath := tests.CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		gitService, err := NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})
}

func TestGetBlameForFileLines(t *testing.T) {
	t.Run("compare terragoat's README second commit", func(t *testing.T) {
		startLine := 0
		endLine := 1
		secondCommitHash := "47accf06f13b503f3bab06fed7860e72f7523cac"

		terragoatPath := tests.CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		gitService, err := NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize repository becauses %s", err)
		}
		generatedGitBlame, err := gitService.GetBlameForFileLines("README.md", []int{startLine, endLine}, secondCommitHash)
		if err != nil {
			t.Errorf("failed to read expected file because %s", err)
		}
		expectedFileLines, err := tests.ReadFileLines("./resources/terragoat_blame_second_commit.txt")
		if err != nil {
			t.Errorf("failed to read expected file because %s", err)
		}

		for lineNum := startLine; lineNum <= endLine; lineNum++ {
			expectedTime := extractUnixDateFromLine(expectedFileLines[lineNum])
			actualTime := generatedGitBlame.BlamesByLine[lineNum].Date
			assert.Equal(t, expectedTime.Unix(), actualTime.Unix())
		}

		if err != nil {
			t.Errorf("could not get latest commit because %s", err)
		}
	})
}
