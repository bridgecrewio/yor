package gitservice

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/pkg/errors"
)

type GitService struct {
	gitRootDir          string
	scanPathFromRoot    string
	repository          *git.Repository
	remoteURL           string
	organization        string
	repoName            string
	BlameByFile         *sync.Map
	PreviousBlameByFile *sync.Map
	currentUserEmail    string
}

var gitGraphLock sync.Mutex

func NewGitService(rootDir string) (*GitService, error) {
	var repository *git.Repository
	var err error
	rootDirIter, _ := filepath.Abs(rootDir)
	for {
		repository, err = git.PlainOpen(rootDirIter)
		if err == nil {
			break
		}
		newRootDir := filepath.Dir(rootDirIter)
		if rootDirIter == newRootDir {
			break
		}
		rootDirIter = newRootDir
	}
	if err != nil {
		return nil, err
	}

	scanAbsDir, _ := filepath.Abs(rootDir)
	scanPathFromRoot, _ := filepath.Rel(rootDirIter, scanAbsDir)

	gitService := GitService{
		gitRootDir:          rootDir,
		scanPathFromRoot:    scanPathFromRoot,
		repository:          repository,
		BlameByFile:         &sync.Map{},
		PreviousBlameByFile: &sync.Map{},
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
		res, _ := filepath.Rel(g.gitRootDir, fp)
		return filepath.Join(g.scanPathFromRoot, res)
	}
	scanPathIter := g.scanPathFromRoot
	parent := filepath.Dir(fp)
	for {
		_, child := filepath.Split(scanPathIter)
		if parent != child {
			break
		}
		scanPathIter, _ = filepath.Split(scanPathIter)
	}
	return filepath.Join(scanPathIter, fp)
}

func (g *GitService) GetBlameForFileLines(filePath string, lines structure.Lines) (*GitBlame, error) {
	logger.Info(fmt.Sprintf("Getting git blame for %v (%v:%v)", filePath, lines.Start, lines.End))
	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	blame, ok := g.BlameByFile.Load(filePath)
	if ok {
		return NewGitBlame(relativeFilePath, filePath, lines, blame.(*git.BlameResult), g), nil
	}

	var err error
	blame, err = g.GetFileBlame(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", filePath, err)
	}

	g.BlameByFile.Store(filePath, blame)

	return NewGitBlame(relativeFilePath, filePath, lines, blame.(*git.BlameResult), g), nil
}

func (g *GitService) GetOrganization() string {
	return g.organization
}

func (g *GitService) GetRepoName() string {
	return g.repoName
}

func wrapGitBlame(selectedCommit *object.Commit, relativeFilePath string) (blame *git.BlameResult, err error) {
	// currently there's a bug inside go-git so in order to mitigate it we wrap it with recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r, string(debug.Stack()))
			err = errors.Errorf("unknown panic, %v", r)
		}
	}()
	blame, err = git.Blame(selectedCommit, relativeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get blame for latest commit of file %s because of error %s", relativeFilePath, err)
	}
	return blame, err
}

func (g *GitService) GetFileBlame(filePath string) (*git.BlameResult, error) {
	blame, ok := g.BlameByFile.Load(filePath)
	if ok {
		return blame.(*git.BlameResult), nil
	}

	relativeFilePath := g.ComputeRelativeFilePath(filePath)
	var selectedCommit *object.Commit

	gitGraphLock.Lock() // Git is a graph, different files can lead to graph scans interfering with each other
	defer gitGraphLock.Unlock()
	head, err := g.repository.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository HEAD for file %s because of error %s", filePath, err)
	}
	selectedCommit, err = g.repository.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to find commit %s ", head.Hash().String())
	}

	parentIter := selectedCommit.Parents()
	previousCommit, err := parentIter.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get previous commit: %s", err)
	}
	blame, err = wrapGitBlame(selectedCommit, relativeFilePath)
	if err != nil {
		return nil, err
	}
	previousBlame, err := wrapGitBlame(previousCommit, relativeFilePath)
	if err != nil {
		return nil, err
	}
	g.BlameByFile.Store(filePath, blame)
	g.PreviousBlameByFile.Store(filePath, previousBlame)

	return blame.(*git.BlameResult), nil
}

func GetGitUserEmail() string {
	log.SetOutput(io.Discard)
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
