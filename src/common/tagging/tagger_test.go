package tagging

import (
	"bridgecrewio/yor/src/common/tagging/tags"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagger(t *testing.T) {
	t.Run("Test tagger skip single tag", func(t *testing.T) {
		tagger := Tagger{SkippedTags: []string{"yor_trace"}}
		tagger.SetTags([]tags.ITag{&tags.Tag{Key: "yor_trace"}, &tags.Tag{Key: "git_modifiers"}})
		tgs := tagger.GetTags()
		assert.Equal(t, 1, len(tgs))
		assert.NotEqual(t, "yor_trace", tgs[0].GetKey())
	})

	t.Run("Test tagger skip regex tags", func(t *testing.T) {
		tagger := Tagger{SkippedTags: []string{"git*"}}
		tagger.SetTags([]tags.ITag{
			&tags.Tag{Key: "yor_trace"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagger.GetTags()
		assert.Equal(t, 1, len(tgs))
		assert.Equal(t, "yor_trace", tgs[0].GetKey())
	})

	t.Run("Test tagger skip multi", func(t *testing.T) {
		tagger := Tagger{SkippedTags: []string{"git*", "yor_trace"}}
		tagger.SetTags([]tags.ITag{
			&tags.Tag{Key: "yor_trace"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagger.GetTags()
		assert.Equal(t, 0, len(tgs))
	})
}
