package structure

import (
	"testing"

	"github.com/bridgecrewio/yor/src/common/tagging/tags"

	"github.com/stretchr/testify/assert"
)

var yorTraceTag = &tags.Tag{
	Key:   "yor_trace",
	Value: "123456789",
}

func TestTagsUpdatedResource(t *testing.T) {
	var existingTags = []tags.ITag{
		&tags.Tag{
			Key:   "git_modifiers",
			Value: "bana",
		},
		&tags.Tag{
			Key:   "git_repo",
			Value: "hatulik",
		},
	}

	traceTag := &tags.Tag{
		Key:   "yor_trace",
		Value: "987654321",
	}
	modifiersTag := &tags.Tag{
		Key:   "git_modifiers",
		Value: "bana/shati",
	}
	var block = Block{
		FilePath:          "/mock.tf",
		ExitingTags:       append(existingTags, yorTraceTag),
		NewTags:           nil,
		RawBlock:          nil,
		IsTaggable:        false,
		TagsAttributeName: "tags",
	}
	t.Run("Test add new tags - skip trace tag", func(t *testing.T) {
		newTags := []tags.ITag{traceTag, modifiersTag}
		block.AddNewTags(newTags)

		assert.Equal(t, 1, len(block.NewTags))
		assert.Equal(t, "123456789", block.GetTraceID())
		block.NewTags = nil
	})

	t.Run("Test add new tags - add trace tag", func(t *testing.T) {
		newTags := []tags.ITag{traceTag, modifiersTag}
		block.AddNewTags(newTags)
		blockTags := block.MergeTags()

		count := 0
		trace := ""
		for _, tag := range blockTags {
			if tag.GetKey() == tags.YorTraceTagKey {
				count++
				trace = tag.GetValue()
			}
		}
		assert.Equal(t, 1, count)
		assert.Equal(t, "123456789", trace)
		block.NewTags = nil
	})

	t.Run("CalculateTagDiff add trace tag", func(t *testing.T) {
		newTags := []tags.ITag{traceTag, modifiersTag}
		block.AddNewTags(newTags)
		tagDiff := block.CalculateTagsDiff()
		assert.Equal(t, 0, len(tagDiff.Added))
		assert.Equal(t, 1, len(tagDiff.Updated))
		assert.Equal(t, "git_modifiers", tagDiff.Updated[0].Key)
		assert.Equal(t, "bana", tagDiff.Updated[0].PrevValue)
		assert.Equal(t, "bana/shati", tagDiff.Updated[0].NewValue)
		block.NewTags = nil
	})
}

func TestTagsNewResource(t *testing.T) {
	var existingTags = []tags.ITag{
		&tags.Tag{
			Key:   "git_modifiers",
			Value: "bana",
		},
		&tags.Tag{
			Key:   "git_repo",
			Value: "hatulik",
		},
	}

	var newTags = []tags.ITag{
		&tags.Tag{
			Key:   "yor_trace",
			Value: "987654321",
		},
		&tags.Tag{
			Key:   "git_modifiers",
			Value: "bana/shati",
		},
	}
	var block = Block{
		FilePath:          "/mock.tf",
		ExitingTags:       existingTags,
		NewTags:           nil,
		RawBlock:          nil,
		IsTaggable:        false,
		TagsAttributeName: "tags",
	}

	t.Run("Test add new tags - add trace tag", func(t *testing.T) {
		block.AddNewTags(newTags)

		assert.Equal(t, 2, len(block.NewTags))
		assert.Equal(t, "987654321", block.GetTraceID())
		block.NewTags = nil
	})

	t.Run("Test add new tags - add trace tag", func(t *testing.T) {
		block.AddNewTags(newTags)
		blockTags := block.MergeTags()

		count := 0
		trace := ""
		for _, tag := range blockTags {
			if tag.GetKey() == tags.YorTraceTagKey {
				count++
				trace = tag.GetValue()
			}
		}
		assert.Equal(t, 1, count)
		assert.Equal(t, "987654321", trace)
		block.NewTags = nil
	})

	t.Run("CalculateTagDiff add trace tag", func(t *testing.T) {
		block.AddNewTags(newTags)
		tagDiff := block.CalculateTagsDiff()
		assert.Equal(t, 1, len(tagDiff.Added))
		assert.Equal(t, 1, len(tagDiff.Updated))
		assert.Equal(t, "yor_trace", tagDiff.Added[0].GetKey())
		assert.Equal(t, "987654321", tagDiff.Added[0].GetValue())
		assert.Equal(t, "git_modifiers", tagDiff.Updated[0].Key)
		assert.Equal(t, "bana", tagDiff.Updated[0].PrevValue)
		assert.Equal(t, "bana/shati", tagDiff.Updated[0].NewValue)
		block.NewTags = nil
	})
}
