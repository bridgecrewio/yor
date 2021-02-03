package git_service

import (
	"bridgecrew.io/yor/git_service"
	"bufio"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const TerragoatUrl = "https://github.com/bridgecrewio/terragoat.git"

func cloneRepo(repoPath string) string {
	dir, err := ioutil.TempDir("", "temp-repo")
	if err != nil {
		log.Fatal(err)
	}

	// Clones the repository into the given dir, just as a normal git clone does
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoPath,
	})

	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func extractUnixDateFromLine(line string) time.Time {
	startIndex := strings.Index(line, "(")
	endIndex := strings.Index(line, ")")
	userDateNum := line[startIndex+1 : endIndex]
	userDateNum = strings.Replace(userDateNum, "  ", " ", -1)
	splitUserDateNum := strings.Split(userDateNum, " ")
	unixString := splitUserDateNum[len(splitUserDateNum)-2]
	intUnix, _ := strconv.ParseInt(unixString, 10, 64)
	unixTime := time.Unix(intUnix, 0)

	return unixTime
}

func TestNewGitService(t *testing.T) {
	t.Run("Get correct organization and repo name", func(t *testing.T) {
		terragoatPath := cloneRepo(TerragoatUrl)

		defer os.RemoveAll(terragoatPath)

		gitService, err := git_service.NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})
}

func TestGetBlameForFileLines(t *testing.T) {
	t.Run("terragoat", func(t *testing.T) {
		terragoatPath := cloneRepo(TerragoatUrl)

		defer os.RemoveAll(terragoatPath)

		gitService, err := git_service.NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize repository becauses %s", err)
		}
		startLine := 0
		endLine := 1
		secondCommitHash := "47accf06f13b503f3bab06fed7860e72f7523cac"
		generatedGitBlame, err := gitService.GetBlameForFileLines("README.md", []int{startLine, endLine}, secondCommitHash)

		expectedBlameFile, err := os.Open("./resources/terragoat_blame_second_commit.txt")
		if err != nil {
			t.Errorf("failed to open expected file because %s", err)
		}

		fileScanner := bufio.NewScanner(expectedBlameFile)
		fileScanner.Split(bufio.ScanLines)
		var fileTextLines []string

		for fileScanner.Scan() {
			fileTextLines = append(fileTextLines, fileScanner.Text())
		}

		expectedBlameFile.Close()

		for lineNum := startLine; lineNum <= endLine; lineNum++ {
			expectedTime := extractUnixDateFromLine(fileTextLines[lineNum])
			actualTime := generatedGitBlame.BlamesByLine[lineNum].Date
			assert.Equal(t, expectedTime.Unix(), actualTime.Unix())
		}

		if err != nil {
			t.Errorf("could not get latest commit because %s", err)
		}
	})
}
