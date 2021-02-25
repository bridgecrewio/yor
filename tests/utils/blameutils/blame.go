package blameutils

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/gitservice"
	"math/rand"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

const Org = "bridgecrewio"
const Repository = "terragoat"
const FilePath = "README.md"
const CommitHash1 = "47accf06f13b503f3bab06fed7860e72f7523cac"
const CommitHash2 = "b2dc884b7439882c4dbe1e660cb1e02a3f84e45d"

func GetGitLines(t *testing.T, numOfLines int) []*git.Line {
	dateStr0 := "2020-03-28T21:42:46.000Z"
	dateStr1 := "2020-03-27T11:56:33.000Z"
	firstCommitDate, err1 := ExtractDate(dateStr0)
	secondCommitDate, err2 := ExtractDate(dateStr1)
	if err1 != nil || err2 != nil {
		assert.Fail(t, "Failed to parse static date")
	}

	results := make([]*git.Line, 0)
	for i := 0; i < numOfLines; i++ {
		if i%3 == 0 {
			results = append(results, &git.Line{
				Author: "schosterbarak@gmail.com",
				Text:   "# Terragoat",
				Date:   firstCommitDate,
				Hash:   plumbing.NewHash(CommitHash1),
			})
		} else {
			results = append(results, &git.Line{
				Author: "jonjozwiak@users.noreply.github.com",
				Text:   "Bridgecrew solution to create vulnerable infrastructure",
				Date:   secondCommitDate,
				Hash:   plumbing.NewHash(CommitHash2),
			})
		}
	}

	return results
}

func SetupBlame(t *testing.T) gitservice.GitBlame {
	gitLines := GetGitLines(t, 3)

	return gitservice.GitBlame{
		GitOrg:        Org,
		GitRepository: Repository,
		FilePath:      FilePath,
		BlamesByLine: map[int]*git.Line{
			0: gitLines[0],
			1: gitLines[1],
			2: gitLines[2],
		},
	}
}

func SetupBlameResults(t *testing.T, path string, numOfLines int) *git.BlameResult {
	return &git.BlameResult{
		Path:  path,
		Rev:   plumbing.Hash{},
		Lines: GetGitLines(t, numOfLines),
	}
}

func ExtractDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	parsedDate, err := time.Parse(layout, dateStr)
	return parsedDate, err
}

func CreateMockBlame(textBytes []byte) git.BlameResult {
	textLines := common.GetLinesFromBytes(textBytes)
	layout = "2006-01-02 15:04:05"
	possibleLines := []*git.Line{
		{
			Author: "shati@gmail.com",
			Date:   getTime("2020-06-16 17:46:24"),
			Hash:   plumbing.NewHash("shati"),
		},
		{
			Author: "bana@gmail.com",
			Date:   getTime("2020-09-25 19:19:02"),
			Hash:   plumbing.NewHash("bana"),
		},
		{
			Author: "checkov@gmail.com",
			Date:   getTime("2020-04-08 19:19:02"),
			Hash:   plumbing.NewHash("checkov"),
		},
	}

	blameLines := make([]*git.Line, 0)
	for _, textLine := range textLines {
		randomIndex := rand.Intn(len(possibleLines))
		selectedLine := possibleLines[randomIndex]
		newLine := git.Line{
			Author: selectedLine.Author,
			Text:   textLine,
			Date:   selectedLine.Date,
			Hash:   selectedLine.Hash,
		}
		blameLines = append(blameLines, &newLine)
	}

	return git.BlameResult{Lines: blameLines}
}

var layout = "2006-01-02 15:04:05"

func getTime(strT string) time.Time {
	t, _ := time.Parse(layout, strT)
	return t
}
