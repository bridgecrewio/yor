package gitservice

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"

	"github.com/go-git/go-git/v5"
)

var CIRegexStrings = []string{
	"\bci\b",
	"[bot]",
	"github-action",
	"\bautomation\b",
}

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
	FilePath      string
	GitUserEmail  string
	CIRegex       *regexp.Regexp
}

func NewGitBlame(filePath string, lines structure.Lines, blameResult *git.BlameResult, gitOrg string, gitRepository string, userEmail string) *GitBlame {
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
	ciRegexp, _ := regexp.Compile("(" + strings.Join(CIRegexStrings, "|") + ")")
	gitBlame.CIRegex = ciRegexp
	return &gitBlame
}

func (g *GitBlame) GetLatestCommit() (latestCommit *git.Line) {
	latestDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, v := range g.BlamesByLine {
		if v == nil {
			// This line was added/edited but not committed yet, so latest commit is nil
			return nil
		}

		isCiBot := g.CIRegex.MatchString(v.Author)
		if latestDate.Before(v.Date) && !isCiBot {
			latestDate = v.Date
			latestCommit = v
		}
	}
	return
}
