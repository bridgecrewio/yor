package gitservice

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitService struct {
	rootDir              string
	repository           *git.Repository
	remoteURL            string
	organization         string
	repoName             string
	blameByFileAndCommit map[string]map[string]*git.BlameResult
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
		rootDir:              rootDir,
		repository:           repository,
		blameByFileAndCommit: make(map[string]map[string]*git.BlameResult),
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

func (g *GitService) GetBlameForFileLines(filePath string, lines []int, commitHash ...string) (*GitBlame, error) {
	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	blame := g.findBlameInCache(relativeFilePath, commitHash)
	if blame != nil {
		return NewGitBlame(relativeFilePath, lines, blame, g.organization, g.repoName), nil
	}

	logOptions := git.LogOptions{
		// order the commits from most to least recent
		Order:    git.LogOrderCommitterTime,
		FileName: &relativeFilePath,
	}
	// fetch commit iterator
	commitIter, err := g.repository.Log(&logOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get log for repository %s because of error %s", g.repository, err)
	}
	var selectedCommit *object.Commit
	if len(commitHash) == 0 {
		// if there no commit was specified, get the latest commit
		selectedCommit, err = commitIter.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest commit for file %s because of error %s", filePath, err)
		}
	} else {
		// find the matching commit
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

	blame = g.findBlameInCache(relativeFilePath, []string{selectedCommit.Hash.String()})
	if blame != nil {
		return NewGitBlame(relativeFilePath, lines, blame, g.organization, g.repoName), nil
	}

	blame, err = git.Blame(selectedCommit, relativeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}
	g.addBlameToCache(relativeFilePath, selectedCommit.Hash.String(), blame)

	return NewGitBlame(relativeFilePath, lines, blame, g.organization, g.repoName), nil

}

func (g *GitService) findBlameInCache(filePath string, commitHash []string) *git.BlameResult {
	if len(commitHash) == 0 {
		return nil
	}
	blame, ok := g.blameByFileAndCommit[filePath][commitHash[0]]
	if !ok {
		return nil
	}

	return blame
}

func (g *GitService) addBlameToCache(filePath string, commitHash string, blame *git.BlameResult) {
	if g.blameByFileAndCommit[filePath] == nil {
		g.blameByFileAndCommit[filePath] = make(map[string]*git.BlameResult)
	}

	g.blameByFileAndCommit[filePath][commitHash] = blame
}

func (g *GitService) GetOrganization() string {
	return g.organization
}

func (g *GitService) GetRepoName() string {
	return g.repoName
}
