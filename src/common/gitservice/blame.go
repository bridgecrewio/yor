package gitservice

import (
	"fmt"
	"time"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"

	"github.com/go-git/go-git/v5"
)

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
	FilePath      string
	GitUserEmail  string
}

func NewGitBlame(filePath string, lines common.Lines, blameResult *git.BlameResult, gitOrg string, gitRepository string, userEmail string) *GitBlame {
	gitBlame := GitBlame{GitOrg: gitOrg, GitRepository: gitRepository, BlamesByLine: map[int]*git.Line{}, FilePath: filePath, GitUserEmail: userEmail}
	startLine := lines.Start - 1 // the lines in blameResult.Lines start from zero while the lines range start from 1
	endLine := lines.End - 1
	for line := startLine; line <= endLine; line++ {
		if line >= len(blameResult.Lines) {
			logger.Warning(fmt.Sprintf("Index out of bound on parsed file %s", filePath))
			return &gitBlame
		}
		gitBlame.BlamesByLine[line+1] = blameResult.Lines[line]
	}

	return &gitBlame
}

func (g *GitBlame) GetLatestCommit() (latestCommit *git.Line) {
	latestDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, v := range g.BlamesByLine {
		if v == nil {
			// This line was added/edited but not committed yet, so latest commit is nil
			return nil
		}
		if latestDate.Before(v.Date) {
			latestDate = v.Date
			latestCommit = v
		}
	}
	return
}
