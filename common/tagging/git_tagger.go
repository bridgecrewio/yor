package tagging

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pmezard/go-difflib/difflib"
)

type GitTagger struct {
	Tagger
	GitService      *gitservice.GitService
	fileLinesMapper map[string]fileLineMapper
}

type fileLineMapper struct {
	originToGit map[int]int
	gitToOrigin map[int]int
}

func (t *GitTagger) InitTagger(path string) {
	gitService, err := gitservice.NewGitService(path)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize git service for path %s", path))
	}
	t.GitService = gitService
}

func (t *GitTagger) initFileMapping(path string) bool {
	fileBlame, err := t.GitService.GetFileBlame(path)
	if err != nil {
		logger.Warning(fmt.Sprintf("Unable to get git blame for file %s: %s", path, err))
		return false
	}

	t.mapOriginFileToGitFile(path, fileBlame)

	return true
}

func (t *GitTagger) CreateTagsForBlock(block structure.IBlock) {
	if _, ok := t.fileLinesMapper[block.GetFilePath()]; !ok {
		t.initFileMapping(block.GetFilePath())
	}
	linesInGit := t.getBlockLinesInGit(block)
	if linesInGit[0] < 0 || linesInGit[1] < 0 {
		return
	}
	blame, err := t.GitService.GetBlameForFileLines(block.GetFilePath(), t.getBlockLinesInGit(block))
	if err != nil {
		logger.Warning(fmt.Sprintf("Failed to tag %v with git tags, err: %v", block.GetResourceID(), err.Error()))
		return
	}
	if blame == nil {
		logger.Warning(fmt.Sprintf("Failed to tag %s with git tags, file must be unstaged", block.GetFilePath()))
		return
	}
	shouldTag := t.updateBlameForOriginLines(block, blame)
	if !shouldTag {
		return
	}

	var newTags []tags.ITag
	for _, tag := range t.Tags {
		newTag, err := tag.CalculateValue(blame)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to calculate tag value of tag %v, err: %s", tag.GetKey(), err))
			continue
		}
		newTags = append(newTags, newTag)
	}
	block.AddNewTags(newTags)
}

func (t *GitTagger) getBlockLinesInGit(block structure.IBlock) []int {
	blockLines := block.GetLines()
	originToGit := t.fileLinesMapper[block.GetFilePath()].originToGit
	originStart := blockLines[0]
	originEnd := blockLines[1]
	gitStart := -1
	gitEnd := -1

	for gitStart == -1 && originStart <= originEnd {
		// find the first mapped line
		gitStart = originToGit[originStart]
		originStart++
	}

	for gitEnd == -1 && originEnd >= blockLines[0] {
		// find the last mapped line
		gitEnd = originToGit[originEnd]
		originEnd--
	}

	return []int{gitStart, gitEnd}
}

// The function maps between the scanned file lines to the lines in the git blame
func (t *GitTagger) mapOriginFileToGitFile(path string, fileBlame *git.BlameResult) {
	if t.fileLinesMapper == nil {
		t.fileLinesMapper = make(map[string]fileLineMapper)
	}
	mapper := fileLineMapper{
		originToGit: make(map[int]int),
		gitToOrigin: make(map[int]int),
	}

	gitLines := make([]string, 0)
	for _, line := range fileBlame.Lines {
		gitLines = append(gitLines, line.Text)
	}

	originFileText, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	originLines := strings.Split(string(originFileText), "\n")

	matcher := difflib.NewMatcher(originLines, gitLines)
	matches := matcher.GetMatchingBlocks()
	currOriginStart := 0
	currGitStart := 0
	for _, match := range matches {
		startInOrigin := match.A
		startInGit := match.B

		// if there were lines that weren't in the range of the previous block, they are changed in the opposite file and will set to -1
		for i := currOriginStart + 1; i <= startInOrigin; i++ {
			mapper.originToGit[i] = -1
		}
		for i := currGitStart + 1; i <= startInGit; i++ {
			mapper.gitToOrigin[i] = -1
		}

		// iterate the matching block and map the corresponding lines
		for i := 1; i <= match.Size; i++ {
			mapper.originToGit[startInOrigin+i] = startInGit + i
			mapper.gitToOrigin[startInGit+i] = startInOrigin + i
		}

		currOriginStart = startInOrigin + match.Size
		currGitStart = startInGit + match.Size
	}

	t.fileLinesMapper[path] = mapper
}

func (t *GitTagger) updateBlameForOriginLines(block structure.IBlock, blame *gitservice.GitBlame) bool {
	gitBlameLines := blame.BlamesByLine
	blockLines := block.GetLines()
	newBlameByLines := make(map[int]*git.Line)
	fileMapping := t.fileLinesMapper[block.GetFilePath()].originToGit

	shouldTag := true
	for blockLine := blockLines[0]; blockLine <= blockLines[1]; blockLine++ {
		if fileMapping[blockLine] == -1 {
			logger.Warning(fmt.Sprintf("unable to tag block in file %s (lines %d-%d) because it contains uncomitted changes", blame.FilePath, blockLines[0], blockLines[1]))
			shouldTag = false
			break
		} else {
			newBlameByLines[blockLine] = gitBlameLines[fileMapping[blockLine]]
		}
	}

	blame.BlamesByLine = newBlameByLines
	return shouldTag
}
