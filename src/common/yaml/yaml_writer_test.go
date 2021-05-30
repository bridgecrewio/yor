package yaml

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation/s3"
	s3tags "github.com/awslabs/goformation/v4/cloudformation/tags"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/simple"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/stretchr/testify/assert"
	"github.com/thepauleh/goserverless/serverless"
)

func TestServerlessWriting(t *testing.T) {
	t.Run("test SLS writing", func(t *testing.T) {
		directory := "../../../tests/serverless/resources/no_tags"
		readFilePath := directory + "/serverless.yml"
		tagGroup := simple.TagGroup{}
		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
		}
		tagGroup.SetTags(extraTags)
		tagGroup.InitTagGroup("", []string{})
		relExpectedPath := directory + "/serverless_tagged.yml"
		slsBlocks := []structure.IBlock{
			&structure.Block{
				FilePath:    readFilePath,
				ExitingTags: nil,
				NewTags:     []tags.ITag{&tags.Tag{Key: "new_tag", Value: "new_value"}},
				RawBlock: serverless.Function{
					Handler: "myFunction.handler",
					Name:    "myFunction",
					Tags: map[string]interface{}{
						"new_tag": "new_value",
					},
				},
				IsTaggable:        true,
				TagsAttributeName: "tags",
				Lines:             structure.Lines{Start: 13, End: 15},
				TagLines:          structure.Lines{Start: -1, End: -1},
			},
			&structure.Block{
				FilePath:    readFilePath,
				ExitingTags: nil,
				NewTags:     []tags.ITag{&tags.Tag{Key: "new_tag", Value: "new_value"}},
				RawBlock: serverless.Function{
					Handler: "myFunction2.handler",
					Name:    "myFunction2",
					Tags: map[string]interface{}{
						"new_tag": "new_value",
					},
				},
				IsTaggable:        true,
				TagsAttributeName: "tags",
				Lines:             structure.Lines{Start: 16, End: 18},
				TagLines:          structure.Lines{Start: -1, End: -1},
			},
		}
		f, _ := ioutil.TempFile(directory, "serverless.*.yaml")
		err := WriteYAMLFile(readFilePath, slsBlocks, f.Name(), structure.Lines{Start: 13, End: 18}, "tags", "functions")
		if err != nil {
			assert.Fail(t, err.Error())
		}
		expectedFilePath, _ := filepath.Abs(relExpectedPath)
		actualFilePath, _ := filepath.Abs(f.Name())
		expected, _ := ioutil.ReadFile(expectedFilePath)
		actualOutput, _ := ioutil.ReadFile(actualFilePath)
		assert.Equal(t, string(expected), string(actualOutput))
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Fail()
			}
		}(f.Name())

	})
}

func TestCFNWriting(t *testing.T) {
	t.Run("test CFN writing", func(t *testing.T) {
		directory := "../../../tests/cloudformation/resources/no_tags"
		readFilePath := directory + "/base.template"
		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
			&tags.Tag{
				Key:   "another_tag",
				Value: "another_val",
			},
		}
		blocks := []structure.IBlock{
			&structure.Block{
				FilePath:    readFilePath,
				ExitingTags: []tags.ITag{},
				NewTags:     extraTags,
				RawBlock: &s3.Bucket{
					AccessControl:                   "PublicRead",
					AWSCloudFormationDeletionPolicy: "Retain",
					WebsiteConfiguration: &s3.Bucket_WebsiteConfiguration{
						ErrorDocument: "error.html",
						IndexDocument: "index.html",
					},
					Tags: []s3tags.Tag{
						{Key: "new_tag", Value: "new_val"},
						{Key: "another_tag", Value: "another_val"},
					},
				},
				IsTaggable:        true,
				TagsAttributeName: "Tags",
				Lines:             structure.Lines{Start: 2, End: 9},
				TagLines:          structure.Lines{Start: -1, End: -1},
				Name:              "S3Bucket",
			},
		}
		f, _ := ioutil.TempFile(directory, "base.*.template")
		err := WriteYAMLFile(readFilePath, blocks, f.Name(), structure.Lines{Start: 2, End: 8}, "Tags", "Resources")
		if err != nil {
			assert.Fail(t, err.Error())
		}
		expectedFilePath := filepath.Join(directory, "expected.txt")
		actualFilePath, _ := filepath.Abs(f.Name())
		expected, _ := ioutil.ReadFile(expectedFilePath)
		actualOutput, _ := ioutil.ReadFile(actualFilePath)
		assert.Equal(t, string(expected), string(actualOutput))
		defer func() {
			_ = os.Remove(f.Name())
		}()

	})
}

func TestExtractIndentationOfLine(t *testing.T) {
	type args struct {
		textLine string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "on indent",
			args: args{textLine: "some text line"},
			want: "",
		},
		{
			name: "3 indents",
			args: args{textLine: "   some text line"},
			want: "   ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractIndentationOfLine(tt.args.textLine); got != tt.want {
				t.Errorf("ExtractIndentationOfLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
