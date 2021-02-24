package gitservice

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
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
		terragoatPath := CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		gitService, err := NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Get correct organization and repo name when in non-root dir", func(t *testing.T) {
		terragoatPath := CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)
		gitService, err := NewGitService(terragoatPath + "/aws")
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Fail if gotten to root dir", func(t *testing.T) {
		terragoatPath := CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		terragoatPath = filepath.Dir(filepath.Dir(terragoatPath))
		gitService, err := NewGitService(terragoatPath)
		assert.NotNil(t, err)
		assert.Nil(t, gitService)
	})
}

func TestGetBlameForFileLines(t *testing.T) {
	t.Run("compare terragoat's README second commit", func(t *testing.T) {
		var err error
		startLine := 1
		endLine := 2
		secondCommitHash := "47accf06f13b503f3bab06fed7860e72f7523cac"

		terragoatPath := CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		gitService, err := NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize repository becauses %s", err)
		}
		generatedGitBlame, err := gitService.GetBlameForFileLines("README.md", []int{startLine, endLine}, secondCommitHash)
		if err != nil {
			t.Errorf("failed to read expected file because %s", err)
		}
		expectedFileLines, err := ReadFileLines("./resources/terragoat_blame_second_commit.txt")
		if err != nil {
			t.Errorf("failed to read expected file because %s", err)
		}

		for lineNum := startLine; lineNum <= endLine; lineNum++ {
			expectedTime := extractUnixDateFromLine(expectedFileLines[lineNum-1])
			actualTime := generatedGitBlame.BlamesByLine[lineNum].Date
			assert.Equal(t, expectedTime.Unix(), actualTime.Unix())
		}

		if err != nil {
			t.Errorf("could not get latest commit because %s", err)
		}
	})
}

func CloneRepo(repoPath string) string {
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

func ReadFileLines(filePath string) ([]string, error) {
	expectedBlameFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(expectedBlameFile)
	fileScanner.Split(bufio.ScanLines)
	fileTextLines := []string{}
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}
	err = expectedBlameFile.Close()
	if err != nil {
		return nil, err
	}

	return fileTextLines, nil
}
