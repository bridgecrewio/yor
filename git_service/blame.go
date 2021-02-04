package git_service

import "github.com/go-git/go-git/v5"

type GitBlame struct {
	GitOrg        string
	GitRepository string
	BlamesByLine  map[int]*git.Line
}

func NewGitBlame(lines []int, blameResult *git.BlameResult, gitOrg string, gitRepository string) *GitBlame {
	// TODO
	return nil
}
