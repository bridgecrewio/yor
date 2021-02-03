package git_service

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"strings"
)

type GitService struct {
	rootDir      string
	repository   *git.Repository
	remoteUrl    string
	organization string
	repoName     string
}

func NewGitService(rootDir string) (*GitService, error) {
	repository, err := git.PlainOpen(rootDir)
	if err != nil {
		return nil, err
	}

	gitService := GitService{
		rootDir:    rootDir,
		repository: repository,
	}

	err = gitService.setOrgAndName()

	return &gitService, err
}

func (g *GitService) setOrgAndName() error {
	remotes, err := g.repository.Remotes()
	if err != nil {
		return fmt.Errorf("failed to get remotes, err: %s", err)
	}

	if len(remotes) > 0 && len(remotes[0].Config().URLs) > 0 {
		g.remoteUrl = remotes[0].Config().URLs[0]
		endpoint, err := transport.NewEndpoint(g.remoteUrl)
		if err != nil {
			return err
		}
		endpointPathParts := strings.Split(strings.TrimSuffix(strings.TrimLeft(endpoint.Path, "/"), ".git"), "/")
		if len(endpointPathParts) != 2 {
			return fmt.Errorf("invalid format of endpoint path: %s", endpoint.Path)
		}
		g.organization = endpointPathParts[0]
		g.repoName = endpointPathParts[1]
	}

	return nil
}

func (g *GitService) GetBlameForFileLines(filePath string, lines []int, commitHash ...string) (*GitBlame, error) {
	logOptions := git.LogOptions{
		Order:    git.LogOrderCommitterTime,
		FileName: &filePath,
	}
	commitIter, err := g.repository.Log(&logOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get log for repository %s because of error %s", g.repository, err)
	}
	var selectedCommit *object.Commit
	if len(commitHash) == 0 {
		selectedCommit, err = commitIter.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest commit for file %s because of error %s", filePath, err)
		}
	} else {
		_ = commitIter.ForEach(func(commit *object.Commit) error {
			if commit.Hash == plumbing.NewHash(commitHash[0]) {
				selectedCommit = commit
			}
			return nil
		})
		if selectedCommit == nil {
			return nil, fmt.Errorf("failed to find commits hash %s in commit logs", commitHash[0])
		}
	}

	blame, err := git.Blame(selectedCommit, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}

	return NewGitBlame(lines, blame, g.organization, g.repoName), nil

}

func (g *GitService) GetOrganization() string {
	return g.organization
}

func (g *GitService) GetRepoName() string {
	return g.repoName
}
