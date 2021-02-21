package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCliArgParsing(t *testing.T) {
	t.Run("Test CLI argument parsing", func(t *testing.T) {
		options := Options{
			Directory:              "",
			Tag:                    "",
			SkipTag:                "",
			CustomTaggers:          nil,
			SkipConfigurationPaths: nil,
			Output:                 "",
			OutputJSONFile:         "",
			ExtraTags:              "",
		}
		assert.Panics(t, options.Validate)
	})

	t.Run("Test CLI argument parsing - valid output", func(t *testing.T) {
		options := Options{
			Directory:              "some/dir",
			Tag:                    "",
			SkipTag:                "",
			CustomTaggers:          nil,
			SkipConfigurationPaths: nil,
			Output:                 "cli",
			OutputJSONFile:         "",
			ExtraTags:              "{}",
		}
		assert.NotPanics(t, options.Validate)
	})

	t.Run("Test CLI argument parsing - invalid output", func(t *testing.T) {
		options := Options{
			Directory:              "some/dir",
			Tag:                    "",
			SkipTag:                "",
			CustomTaggers:          nil,
			SkipConfigurationPaths: nil,
			Output:                 "junitxml",
			OutputJSONFile:         "",
			ExtraTags:              "",
		}
		assert.Panics(t, options.Validate)
	})
}
