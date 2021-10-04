package structure

import (
	goformationTags "github.com/awslabs/goformation/v5/cloudformation/tags"
	"sort"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"

	"github.com/awslabs/goformation/v5"
	"github.com/stretchr/testify/assert"
)

func TestCloudformationBlock_MergeCFNTags(t *testing.T) {
	t.Run("test merging with the same order", func(t *testing.T) {
		existingTags := []tags.ITag{
			&tags.Tag{Key: "cfn_tag_1", Value: "1"}, &tags.Tag{Key: "cfn_tag_2", Value: "2"}, &tags.Tag{Key: "yor_trace", Value: "should not change"}, &tags.Tag{Key: "git_last_modified_at", Value: "1"},
		}
		newTags := []tags.ITag{
			&tags.Tag{Key: "yor_trace", Value: "2"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
		}

		expectedMergedTags := []tags.ITag{
			&tags.Tag{Key: "cfn_tag_1", Value: "1"}, &tags.Tag{Key: "cfn_tag_2", Value: "2"}, &tags.Tag{Key: "yor_trace", Value: "should not change"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
		}
		b := &CloudformationBlock{
			Block: structure.Block{
				ExitingTags: existingTags,
				NewTags:     newTags,
			},
		}
		actualMergedTags := b.MergeTags()
		for i, expectedTag := range expectedMergedTags {
			assert.Equal(t, expectedTag.GetKey(), actualMergedTags[i].GetKey())
			assert.Equal(t, expectedTag.GetValue(), actualMergedTags[i].GetValue())
		}
	})
}

func updateCFNTags(templatePath string, t *testing.T) {
	existingTags := []tags.ITag{
		&tags.Tag{Key: "MyTag", Value: "TagValue"},
	}
	newTags := []tags.ITag{
		&tags.Tag{Key: "yor_trace", Value: "yor_trace"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
	}

	expectedMergedTags := []tags.ITag{
		&tags.Tag{Key: "MyTag", Value: "TagValue"}, &tags.Tag{Key: "yor_trace", Value: "yor_trace"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
	}

	template, err := goformation.Open(templatePath)
	if err != nil {
		t.Errorf("There was an error processing the cloudformation template: %s", err)
	}
	resourceName := "NewVolume"
	resource := template.Resources[resourceName]

	b := &CloudformationBlock{
		Block: structure.Block{
			ExitingTags:       existingTags,
			NewTags:           newTags,
			RawBlock:          resource,
			IsTaggable:        true,
			TagsAttributeName: "Tags",
			Lines:             structure.Lines{Start: 4, End: 14},
			Name:              resourceName,
		},
	}

	b.UpdateTags()

	currentRawBlock := b.RawBlock.(map[string]interface{})
	currentTags := currentRawBlock["Properties"].(map[string]interface{})["Tags"].([]goformationTags.Tag)
	sort.Slice(expectedMergedTags, func(i, j int) bool {
		return expectedMergedTags[i].GetKey() > expectedMergedTags[j].GetKey()
	})

	sort.Slice(currentTags, func(i, j int) bool {
		return currentTags[i].Key > currentTags[j].Key
	})

	assert.Equal(t, len(expectedMergedTags), len(currentTags))
	for i, expectedTag := range expectedMergedTags {
		assert.Equal(t, expectedTag.GetKey(), currentTags[i].Key)
		assert.Equal(t, expectedTag.GetValue(), currentTags[i].Value)
	}
}

func TestCloudformationBlock_UpdateTags(t *testing.T) {

	t.Run("update cfn tags", func(t *testing.T) {
		updateCFNTags("../../../tests/cloudformation/resources/ebs/ebs.yaml", t)
	})
	t.Run("update cfn json tags", func(t *testing.T) {
		updateCFNTags("../../../tests/cloudformation/resources/ebs/ebs.json", t)
	})
}
