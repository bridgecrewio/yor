package gitservice

import (
	"bridgecrewio/yor/common/logger"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitService struct {
	rootDir      string
	repository   *git.Repository
	remoteURL    string
	organization string
	repoName     string
	BlameByFile  map[string]*git.BlameResult
}

func NewGitService(rootDir string) (*GitService, error) {
	var repository *git.Repository
	var err error
	for {
		repository, err = git.PlainOpen(rootDir)
		if err == nil {
			break
		}

		newRootDir := filepath.Dir(rootDir)
		if rootDir == newRootDir {
			break
		}
		rootDir = newRootDir
	}
	if err != nil {
		return nil, err
	}

	gitService := GitService{
		rootDir:     rootDir,
		repository:  repository,
		BlameByFile: make(map[string]*git.BlameResult),
	}

	err = gitService.setOrgAndName()

	return &gitService, err
}

func (g *GitService) setOrgAndName() error {
	// get remotes to find the repository's url
	remotes, err := g.repository.Remotes()
	if err != nil {
		return fmt.Errorf("failed to get remotes, err: %s", err)
	}

	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			g.remoteURL = remote.Config().URLs[0]
			// get endpoint structured like '/bridgecrewio/yor.git
			endpoint, err := transport.NewEndpoint(g.remoteURL)
			if err != nil {
				return err
			}
			// remove leading '/' from path and trailing '.git. suffix, then split by '/'
			endpointPathParts := strings.Split(strings.TrimSuffix(strings.TrimLeft(endpoint.Path, "/"), ".git"), "/")
			if len(endpointPathParts) != 2 {
				return fmt.Errorf("invalid format of endpoint path: %s", endpoint.Path)
			}
			g.organization = endpointPathParts[0]
			g.repoName = endpointPathParts[1]
			break
		}
	}

	return nil
}
func (g *GitService) ComputeRelativeFilePath(filepath string) string {
	return strings.ReplaceAll(filepath, fmt.Sprintf("%s/", g.rootDir), "")
}

func (g *GitService) GetBlameForFileLines(filePath string, lines []int) (*GitBlame, error) {
	logger.Info(fmt.Sprintf("Getting git blame for %v (%v:%v)", filePath, lines[0], lines[1]))
	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	blame, ok := g.BlameByFile[filePath]
	if ok {
		return NewGitBlame(relativeFilePath, lines, blame, g.organization, g.repoName), nil
	}

	var err error
	blame, err = g.GetFileBlame(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}

	g.BlameByFile[filePath] = blame

	return NewGitBlame(relativeFilePath, lines, blame, g.organization, g.repoName), nil
}

func (g *GitService) GetOrganization() string {
	return g.organization
}

func (g *GitService) GetRepoName() string {
	return g.repoName
}

func (g *GitService) GetFileBlame(filePath string) (*git.BlameResult, error) {
	blame, ok := g.BlameByFile[filePath]
	if ok {
		return blame, nil
	}

	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	var selectedCommit *object.Commit

	head, err := g.repository.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository HEAD for file %s because of error %s", filePath, err)
	}
	selectedCommit, err = g.repository.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to find commit %s ", head.Hash().String())
	}

	blame, err = git.Blame(selectedCommit, relativeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}
	g.BlameByFile[filePath] = blame

	return blame, nil
}
