package gitservice

import (
	"bridgecrewio/yor/common/logger"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
)

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
	FilePath      string
}

// lines: []int{startLine, endLine}
func NewGitBlame(filePath string, lines []int, blameResult *git.BlameResult, gitOrg string, gitRepository string) *GitBlame {
	gitBlame := GitBlame{GitOrg: gitOrg, GitRepository: gitRepository, BlamesByLine: map[int]*git.Line{}, FilePath: filePath}
	i := 0
	for line := lines[0]; line <= lines[1]; line++ {
		if i < 0 || i >= len(blameResult.Lines) {
			logger.Error(fmt.Sprintf("Index out of bound on parsed file %s", filePath))
		}
		gitBlame.BlamesByLine[line] = blameResult.Lines[i]
		i++
	}

	return &gitBlame
}

func (g *GitBlame) GetLatestCommit() (latestCommit *git.Line) {
	latestDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, v := range g.BlamesByLine {
		if latestDate.Before(v.Date) {
			latestDate = v.Date
			latestCommit = v
		}
	}
	return
}
