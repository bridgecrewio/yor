package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerlessBlock_MergeSLSTags(t *testing.T) {
	t.Run("test merging with the same order", func(t *testing.T) {
		existingTags := []tags.ITag{
			&tags.Tag{Key: "sls_tag_1", Value: "1"}, &tags.Tag{Key: "sls_tag_2", Value: "2"}, &tags.Tag{Key: "yor_trace", Value: "should not change"}, &tags.Tag{Key: "git_last_modified_at", Value: "1"},
		}
		newTags := []tags.ITag{
			&tags.Tag{Key: "yor_trace", Value: "2"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
		}

		expectedMergedTags := []tags.ITag{
			&tags.Tag{Key: "sls_tag_1", Value: "1"}, &tags.Tag{Key: "sls_tag_2", Value: "2"}, &tags.Tag{Key: "yor_trace", Value: "should not change"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
		}
		b := &ServerlessBlock{
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

func TestServerlessBlock_UpdateTags(t *testing.T) {
	t.Run("update sls tags", func(t *testing.T) {
		parser := ServerlessParser{}
		parser.Init("../../tests/serverless/resources", nil)
		existingTags := []tags.ITag{
			&tags.Tag{Key: "TAG1_FUNC", Value: "Func1 Tag Value"},
			&tags.Tag{Key: "TAG2_FUNC", Value: "Func2 Tag Value"},
		}
		newTags := []tags.ITag{
			&tags.Tag{Key: "yor_trace", Value: "yor_trace"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
		}

		expectedMergedTags := []tags.ITag{
			&tags.Tag{Key: "TAG1_FUNC", Value: "Func1 Tag Value"}, &tags.Tag{Key: "TAG2_FUNC", Value: "Func2 Tag Value"}, &tags.Tag{Key: "yor_trace", Value: "yor_trace"}, &tags.Tag{Key: "git_last_modified_at", Value: "2"},
		}
		fmt.Println(expectedMergedTags)
		absFilePath, _ := filepath.Abs(strings.Join([]string{parser.rootDir, "serverless.yml"}, "/"))
		template, err := parser.ParseFile(absFilePath)
		if err != nil {
			t.Errorf("There was an error processing the cloudformation template: %s", err)
		}
		template = append(template, nil)
		resourceName := "myFunction"
		resource := template

		b := &ServerlessBlock{
			Block: structure.Block{
				ExitingTags:       existingTags,
				NewTags:           newTags,
				RawBlock:          resource[0].GetRawBlock(),
				IsTaggable:        true,
				TagsAttributeName: "tags",
			},
			name:  resourceName,
			lines: common.Lines{Start: 4, End: 14},
		}

		b.UpdateTags()

		//currentRawBlock := b.RawBlock.(*ec2.Volume)
		//currentTags := currentRawBlock.Tags
		//sort.Slice(expectedMergedTags, func(i, j int) bool {
		//	return expectedMergedTags[i].GetKey() > expectedMergedTags[j].GetKey()
		//})
		//
		//sort.Slice(currentTags, func(i, j int) bool {
		//	return currentTags[i].Key > currentTags[j].Key
		//})
		//
		//assert.Equal(t, len(expectedMergedTags), len(currentTags))
		//for i, expectedTag := range expectedMergedTags {
		//	assert.Equal(t, expectedTag.GetKey(), currentTags[i].Key)
		//	assert.Equal(t, expectedTag.GetValue(), currentTags[i].Value)
		//}

	})
}
