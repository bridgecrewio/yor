package gitservice

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitService struct {
	gitRootDir       string
	scanPathFromRoot string
	repository       *git.Repository
	remoteURL        string
	organization     string
	repoName         string
	BlameByFile      map[string]*git.BlameResult
	currentUserEmail string
}

func NewGitService(rootDir string) (*GitService, error) {
	var repository *git.Repository
	var err error
	scanPathFromRoot := "."
	for {
		repository, err = git.PlainOpen(rootDir)
		if err == nil {
			break
		}
		absRoot, _ := filepath.Abs(rootDir)
		newRootDir := filepath.Dir(absRoot)
		relPath, _ := filepath.Rel(newRootDir, absRoot)
		scanPathFromRoot = filepath.Join(relPath, scanPathFromRoot)
		if rootDir == newRootDir {
			break
		}
		rootDir = newRootDir
	}
	if err != nil {
		return nil, err
	}

	gitService := GitService{
		gitRootDir:       rootDir,
		scanPathFromRoot: scanPathFromRoot,
		repository:       repository,
		BlameByFile:      make(map[string]*git.BlameResult),
	}

	err = gitService.setOrgAndName()
	gitService.currentUserEmail = GetGitUserEmail()

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
			// get endpoint structured like '/github.com/bridgecrewio/yor.git
			endpoint, err := transport.NewEndpoint(g.remoteURL)
			if err != nil {
				return err
			}
			// remove leading '/' from path and trailing '.git. suffix, then split by '/'
			endpointPathParts := strings.Split(strings.TrimSuffix(strings.TrimLeft(endpoint.Path, "/"), ".git"), "/")
			if len(endpointPathParts) < 2 {
				return fmt.Errorf("invalid format of endpoint path: %s", endpoint.Path)
			}
			g.organization = endpointPathParts[0]
			g.repoName = strings.Join(endpointPathParts[1:], "/")
			break
		}
	}

	return nil
}

func (g *GitService) ComputeRelativeFilePath(fp string) string {
	if strings.HasPrefix(fp, g.gitRootDir) {
		return strings.ReplaceAll(fp, fmt.Sprintf("%s/", g.gitRootDir), "")
	}
	return filepath.Join(g.scanPathFromRoot, fp)
}

func (g *GitService) GetBlameForFileLines(filePath string, lines structure.Lines) (*GitBlame, error) {
	logger.Info(fmt.Sprintf("Getting git blame for %v (%v:%v)", filePath, lines.Start, lines.End))
	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	blame, ok := g.BlameByFile[filePath]
	if ok {
		return NewGitBlame(relativeFilePath, lines, blame, g.organization, g.repoName, g.currentUserEmail), nil
	}

	var err error
	blame, err = g.GetFileBlame(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}

	g.BlameByFile[filePath] = blame

	return NewGitBlame(relativeFilePath, lines, blame, g.organization, g.repoName, g.currentUserEmail), nil
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

func GetGitUserEmail() string {
	log.SetOutput(ioutil.Discard)
	cmd := exec.Command("git", "config", "user.email")
	email, err := cmd.Output()
	stdout := os.Stdout
	log.SetOutput(stdout)
	if err != nil {
		logger.Debug(fmt.Sprintf("unable to get current git user email: %s", err))
		return ""
	}
	return strings.ReplaceAll(string(email), "\n", "")
}
