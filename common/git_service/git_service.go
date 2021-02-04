package git_service

import (
	"github.com/go-git/go-git/v5"
)

type GitService struct {
	rootDir      string
	repository   *git.Repository
	remoteUrl    string
	organization string
	repoName     string
}

func NewGitService(rootDir string) (*GitService, error) {
	// TODO
	return nil, nil
}

func (g *GitService) GetBlameForFileLines(filePath string, lines []int, commitHash ...string) (*GitBlame, error) {
	// TODO
	return nil, nil
}
