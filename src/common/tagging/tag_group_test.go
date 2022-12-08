package tagging

import (
	"testing"

	"github.com/bridgecrewio/yor/src/common/tagging/tags"

	"github.com/stretchr/testify/assert"
)

func TestTagGroup(t *testing.T) {
	t.Run("Test tagGroup skip single tag", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"yor_trace"}}
		tagGroup.SetTags([]tags.ITag{&tags.Tag{Key: "yor_trace"}, &tags.Tag{Key: "git_modifiers"}})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 1, len(tgs))
		assert.NotEqual(t, "yor_trace", tgs[0].GetKey())
	})

	t.Run("Test tagGroup skip regex tags", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"git*"}}
		tagGroup.SetTags([]tags.ITag{
			&tags.Tag{Key: "yor_trace"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 1, len(tgs))
		assert.Equal(t, "yor_trace", tgs[0].GetKey())
	})

	t.Run("Test tagGroup skip multi", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"git*", "yor_trace"}}
		tagGroup.SetTags([]tags.ITag{
			&tags.Tag{Key: "yor_trace"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 0, len(tgs))
	})

	t.Run("Test tag prefix not broke tagGroup skip multi", func(t *testing.T) {
		tagGroup := TagGroup{SkippedTags: []string{"git*", "yor_trace"}, Options: InitTagGroupOptions{
			TagPrefix: "prefix_",
		}}
		tagGroup.SetTags([]tags.ITag{
			&tags.Tag{Key: "yor_trace"},
			&tags.Tag{Key: "git_modifiers"},
			&tags.Tag{Key: "git_modifiers"},
		})
		tgs := tagGroup.GetTags()
		assert.Equal(t, 0, len(tgs))
	})
}
