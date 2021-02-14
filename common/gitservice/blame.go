package gitservice

import "github.com/go-git/go-git/v5"

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
	FilePath      string
}

// lines: []int{startLine, endLine}
func NewGitBlame(filePath string, lines []int, blameResult *git.BlameResult, gitOrg string, gitRepository string) *GitBlame {
	gitBlame := GitBlame{GitOrg: gitOrg, GitRepository: gitRepository, BlamesByLine: map[int]*git.Line{}, FilePath: filePath}
	for line := lines[0]; line <= lines[1]; line++ {
		gitBlame.BlamesByLine[line] = blameResult.Lines[line]
	}

	return &gitBlame
}
