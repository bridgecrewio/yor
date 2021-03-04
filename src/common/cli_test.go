package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliArgParsing(t *testing.T) {
	t.Run("Test tag argument parsing - valid output", func(t *testing.T) {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            "",
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "cli",
			OutputJSONFile: "",
		}
		assert.NotPanics(t, options.Validate)
	})

	t.Run("Test tag argument parsing - invalid output", func(t *testing.T) {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            "",
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "junitxml",
			OutputJSONFile: "",
		}
		assert.Panics(t, options.Validate)
	})

	t.Run("Test tag argument parsing - valid tag groups", func(t *testing.T) {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            "",
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "cli",
			OutputJSONFile: "",
			TagGroups:      []string{"git", "code2cloud"},
		}
		assert.NotPanics(t, options.Validate)
	})

	t.Run("Test tag argument parsing - invalid tag groups", func(t *testing.T) {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            "",
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "junitxml",
			OutputJSONFile: "",
			TagGroups:      []string{"git", "custom"},
		}
		assert.Panics(t, options.Validate)
	})

	t.Run("Test CLI argument parsing - list-tags - invalid output", func(t *testing.T) {
		options := ListTagsOptions{
			TagGroups: []string{"custom"},
		}
		assert.Panics(t, options.Validate)
	})

	t.Run("Test CLI argument parsing - list-tags - valid output", func(t *testing.T) {
		options := ListTagsOptions{
			TagGroups: []string{"simple", "git"},
		}
		assert.NotPanics(t, options.Validate)
	})
}
