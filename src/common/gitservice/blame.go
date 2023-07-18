package gitservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"

	"github.com/go-git/go-git/v5"
)

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
	FilePath      string
	GitUserEmail  string
}

func GetPreviousBlameResult(g *GitService, filePath string) (*git.BlameResult, *object.Commit) {
	if g.repository == nil {
		return nil, nil
	}
	ref, err := g.repository.Head()
	if err != nil {
		return nil, nil
	}
	commit, err := g.repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, nil
	}
	parentIter := commit.Parents()
	previousCommit, err := parentIter.Next()
	if err != nil {
		return nil, nil
	}

	var previousBlameResult *git.BlameResult
	result, _ := g.PreviousBlameByFile.Load(filePath)
	previousBlameResult = result.(*git.BlameResult)
	return previousBlameResult, previousCommit
}

func NewGitBlame(relativeFilePath string, filePath string, lines structure.Lines, blameResult *git.BlameResult, g *GitService) *GitBlame {
	gitBlame := GitBlame{GitOrg: g.organization, GitRepository: g.repoName, BlamesByLine: map[int]*git.Line{}, FilePath: relativeFilePath, GitUserEmail: g.currentUserEmail}
	startLine := lines.Start - 1 // the lines in blameResult.Lines start from zero while the lines range start from 1
	endLine := lines.End - 1
	previousBlameResult, previousCommit := GetPreviousBlameResult(g, filePath)

	for line := startLine; line <= endLine; line++ {
		if line >= len(blameResult.Lines) {
			logger.Warning(fmt.Sprintf("Index out of bound on parsed file %s", relativeFilePath))
			return &gitBlame
		}
		gitBlame.BlamesByLine[line+1] = blameResult.Lines[line]

		// Check if the line has been removed in the current state of the file
		if previousBlameResult != nil && len(previousBlameResult.Lines) > len(blameResult.Lines) {
			if previousBlameResult.Lines[line].Text != blameResult.Lines[line].Text {
				// The line has been removed, so update the git commit id
				gitBlame.BlamesByLine[line+1].Hash = previousCommit.Hash
			}
		}
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
		if latestDate.Before(v.Date) &&
			// Commit was not made by CI, i.e. github actions (for now)
			!strings.Contains(v.Author, "[bot]") && !strings.Contains(v.Author, "github-actions") {
			latestDate = v.Date
			latestCommit = v
		}
	}
	return
}
