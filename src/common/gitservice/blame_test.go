package gitservice

import (
	"sync"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

// newMinimalGitService builds a GitService with just the fields NewGitBlame and
// GetPreviousBlameResult touch, so we can unit-test the line-range handling
// without cloning a real repository.
func newMinimalGitService() *GitService {
	return &GitService{
		repository:          nil, // GetPreviousBlameResult is nil-safe on a nil repository
		organization:        "test-org",
		repoName:            "test-repo",
		BlameByFile:         &sync.Map{},
		PreviousBlameByFile: &sync.Map{},
		currentUserEmail:    "test@example.com",
	}
}

func newBlameResultWithLines(n int) *git.BlameResult {
	lines := make([]*git.Line, 0, n)
	for i := 0; i < n; i++ {
		lines = append(lines, &git.Line{Text: "line", Author: "a@b.c"})
	}
	return &git.BlameResult{Lines: lines}
}

// TestNewGitBlameDegenerateLineRanges is a regression test for a negative-index
// panic ("index out of range [-1]") that occurred when a block's lines could
// not be mapped to the git blame (e.g. a Helm document with no locatable
// metadata), leaving Start == 0 and producing startLine == -1.
func TestNewGitBlameDegenerateLineRanges(t *testing.T) {
	gitSvc := newMinimalGitService()
	blameResult := newBlameResultWithLines(5)

	t.Run("Start of 0 does not panic and yields no blames", func(t *testing.T) {
		var gitBlame *GitBlame
		assert.NotPanics(t, func() {
			gitBlame = NewGitBlame("a.yaml", "/tmp/a.yaml", structure.Lines{Start: 0, End: 3}, blameResult, gitSvc)
		})
		assert.NotNil(t, gitBlame)
		assert.Empty(t, gitBlame.BlamesByLine)
	})

	t.Run("Negative start does not panic and yields no blames", func(t *testing.T) {
		var gitBlame *GitBlame
		assert.NotPanics(t, func() {
			gitBlame = NewGitBlame("a.yaml", "/tmp/a.yaml", structure.Lines{Start: -2, End: 1}, blameResult, gitSvc)
		})
		assert.NotNil(t, gitBlame)
		assert.Empty(t, gitBlame.BlamesByLine)
	})

	t.Run("End beyond blame length does not panic", func(t *testing.T) {
		var gitBlame *GitBlame
		assert.NotPanics(t, func() {
			gitBlame = NewGitBlame("a.yaml", "/tmp/a.yaml", structure.Lines{Start: 1, End: 100}, blameResult, gitSvc)
		})
		assert.NotNil(t, gitBlame)
	})

	t.Run("Valid range maps the expected lines", func(t *testing.T) {
		gitBlame := NewGitBlame("a.yaml", "/tmp/a.yaml", structure.Lines{Start: 1, End: 3}, blameResult, gitSvc)
		assert.NotNil(t, gitBlame)
		assert.Len(t, gitBlame.BlamesByLine, 3)
		assert.Contains(t, gitBlame.BlamesByLine, 1)
		assert.Contains(t, gitBlame.BlamesByLine, 3)
	})
}
