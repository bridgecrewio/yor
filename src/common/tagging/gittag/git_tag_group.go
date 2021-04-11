package gittag

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5"
	"github.com/pmezard/go-difflib/difflib"
)

type TagGroup struct {
	tagging.TagGroup
	GitService      *gitservice.GitService
	fileLinesMapper map[string]fileLineMapper
}

type fileLineMapper struct {
	originToGit map[int]int
	gitToOrigin map[int]int
}

func (t *TagGroup) InitTagGroup(path string, skippedTags []string) {
	t.SkippedTags = skippedTags
	gitService, err := gitservice.NewGitService(path)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize git service for path %s", path))
	}
	t.GitService = gitService
	t.SetTags(t.GetDefaultTags())
}

func (t *TagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{
		&GitOrgTag{},
		&GitRepoTag{},
		&GitFileTag{},
		&GitCommitTag{},
		&GitModifiersTag{},
		&GitLastModifiedAtTag{},
		&GitLastModifiedByTag{},
	}
}

func (t *TagGroup) initFileMapping(path string) bool {
	fileBlame, err := t.GitService.GetFileBlame(path)
	if err != nil {
		logger.Warning(fmt.Sprintf("Unable to get git blame for file %s: %s", path, err))
		return false
	}

	t.mapOriginFileToGitFile(path, fileBlame)

	return true
}

func (t *TagGroup) CreateTagsForBlock(block structure.IBlock) {
	if _, ok := t.fileLinesMapper[block.GetFilePath()]; !ok {
		t.initFileMapping(block.GetFilePath())
	}
	linesInGit := t.getBlockLinesInGit(block)
	if linesInGit.Start < 0 || linesInGit.End < 0 {
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
	t.updateBlameForOriginLines(block, blame)
	if !t.hasNonYorChanges(blame, block) {
		return
	}
	var newTags []tags.ITag
	for _, tag := range t.GetTags() {
		newTag, err := tag.CalculateValue(blame)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to calculate tag value of tag %v, err: %s", tag.GetKey(), err))
			continue
		}
		newTags = append(newTags, newTag)
	}
	block.AddNewTags(newTags)
}

func (t *TagGroup) getBlockLinesInGit(block structure.IBlock) common.Lines {
	blockLines := block.GetLines()
	originToGit := t.fileLinesMapper[block.GetFilePath()].originToGit
	originStart := blockLines.Start
	originEnd := blockLines.End
	gitStart := -1
	gitEnd := -1

	for gitStart == -1 && originStart <= originEnd {
		// find the first mapped line
		gitStart = originToGit[originStart]
		originStart++
	}

	for gitEnd == -1 && originEnd >= blockLines.Start {
		// find the last mapped line
		gitEnd = originToGit[originEnd]
		originEnd--
	}

	return common.Lines{Start: gitStart, End: gitEnd}
}

// The function maps between the scanned file lines to the lines in the git blame
func (t *TagGroup) mapOriginFileToGitFile(path string, fileBlame *git.BlameResult) {
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

	originLines := common.GetLinesFromBytes(originFileText)

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

func (t *TagGroup) updateBlameForOriginLines(block structure.IBlock, blame *gitservice.GitBlame) {
	gitBlameLines := blame.BlamesByLine
	blockLines := block.GetLines(true)
	newBlameByLines := make(map[int]*git.Line)
	fileMapping := t.fileLinesMapper[block.GetFilePath()].originToGit

	for blockLine := blockLines.Start; blockLine <= blockLines.End; blockLine++ {
		if fileMapping[blockLine] == -1 {
			newBlameByLines[blockLine] = &git.Line{
				Author: blame.GitUserEmail,
				Date:   time.Now().UTC(),
				Hash:   plumbing.ZeroHash,
			}
		} else {
			newBlameByLines[blockLine] = gitBlameLines[fileMapping[blockLine]]
		}
	}

	blame.BlamesByLine = newBlameByLines
}

func (t *TagGroup) hasNonYorChanges(blame *gitservice.GitBlame, block structure.IBlock) bool {
	allTagsKeysStr := os.Getenv("TAG_KEYS")
	allTagsKeys := strings.Split(allTagsKeysStr, ",")
	tagsLines := block.GetTagsLines()
	for lineNum, line := range blame.BlamesByLine {
		if line.Hash.String() != blame.LatestCommit {
			continue
		}
		if lineNum < tagsLines.Start || lineNum > tagsLines.End {
			return true
		}
		existingTags := block.GetExistingTags()
		for _, tag := range existingTags {
			for _, linePart := range strings.Split(line.Text, block.GetSeparator()) {
				trimmedLinePart := strings.TrimSpace(linePart)
				if trimmedLinePart == tag.GetKey() || trimmedLinePart == tag.GetValue() {
					if !common.InSlice(allTagsKeys, tag.GetKey()) {
						return true
					}
					break
				}
			}
		}
	}

	return false
}
