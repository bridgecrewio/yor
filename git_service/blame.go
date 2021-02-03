package git_service

import "github.com/go-git/go-git/v5"

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
}

func NewGitBlame(lines []int, blameResult *git.BlameResult, gitOrg string, gitRepository string) *GitBlame {
	gitBlame := GitBlame{GitOrg: gitOrg, GitRepository: gitRepository, BlamesByLine: map[int]*git.Line{}}
	for line := lines[0]; line <= lines[1]; line++ {
		gitBlame.BlamesByLine[line] = blameResult.Lines[line]
	}

	return &gitBlame
}
