package clioptions

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliArgParsing(t *testing.T) {
	t.Run("Tag local options flag", func(t *testing.T) {
		options := TagOptions{
			Directory:       "some/dir",
			Tag:             nil,
			SkipTags:        nil,
			CustomTagging:   nil,
			SkipDirs:        nil,
			Output:          "cli",
			OutputJSONFile:  "",
			ConfigFile:      "",
			DryRun:          true,
			TagLocalModules: true,
			TagPrefix:       "prefix",
		}
		// Expect the validation to pass without throwing errors
		options.Validate()
	})
	t.Run("Test tag argument parsing - valid output", func(t *testing.T) {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            nil,
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "cli",
			OutputJSONFile: "",
			ConfigFile:     "",
			DryRun:         true,
			TagPrefix:      "",
		}
		// Expect the validation to pass without throwing errors
		options.Validate()
	})

	t.Run("Test tag argument parsing - invalid output", func(t *testing.T) {
		cmd := exec.Command(os.Args[0], "-test.run=TestOutputCrasher")
		cmd.Env = append(cmd.Env, "UT_CRASH=RUN")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		assert.Fail(t, "Should have failed already")
	})

	t.Run("Test tag argument parsing - valid tag groups", func(t *testing.T) {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            nil,
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "cli",
			OutputJSONFile: "",
			TagGroups:      []string{"git", "code2cloud"},
		}
		// Expect the validation to pass without throwing errors
		options.Validate()
	})

	t.Run("Test tag argument parsing - invalid tag groups", func(t *testing.T) {
		cmd := exec.Command(os.Args[0], "-test.run=TestOutputCrasher")
		cmd.Env = append(cmd.Env, "UT_CRASH=RUN")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		assert.Fail(t, "Should have failed already")
	})

	t.Run("Test CLI argument parsing - list-tags - invalid output", func(t *testing.T) {
		cmd := exec.Command(os.Args[0], "-test.run=TestOutputCrasher")
		cmd.Env = append(cmd.Env, "UT_CRASH=RUN")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		assert.Fail(t, "Should have failed already")
	})

	t.Run("Test CLI argument parsing - list-tags - valid output", func(t *testing.T) {
		options := ListTagsOptions{
			TagGroups: []string{"simple", "git"},
		}
		// Expect the validation to pass without throwing errors
		options.Validate()
	})
}

func TestOutputCrasher(t *testing.T) {
	if os.Getenv("UT_CRASH") == "RUN" {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            nil,
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "junitxml",
			OutputJSONFile: "",
			TagGroups:      []string{"git", "custom"},
			DryRun:         true,
		}
		options.Validate()
	}
}

func TestTagGroupCrasher(t *testing.T) {
	if os.Getenv("UT_CRASH") == "RUN" {
		options := TagOptions{
			Directory:      "some/dir",
			Tag:            nil,
			SkipTags:       nil,
			CustomTagging:  nil,
			SkipDirs:       nil,
			Output:         "cli",
			OutputJSONFile: "",
			TagGroups:      []string{"git", "custom"},
			DryRun:         true,
		}
		options.Validate()
	}
}

func TestListTagsGroupCrasher(t *testing.T) {
	if os.Getenv("UT_CRASH") == "RUN" {
		options := ListTagsOptions{
			TagGroups: []string{"custom"},
		}
		options.Validate()
	}
}
