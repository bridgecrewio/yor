package json

import (
	"github.com/bridgecrewio/yor/src/common/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/simple"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/stretchr/testify/assert"
	"github.com/thepauleh/goserverless/serverless"
)

func TestCloudformationWriting(t *testing.T) {
	t.Run("test CFN writing", func(t *testing.T) {
		directory := "../../../tests/cloudformation/resources/ec2_untagged"
		readFilePath := directory + "/ec2_untagged.json"
		tagGroup := simple.TagGroup{}
		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
		}
		tagGroup.SetTags(extraTags)
		tagGroup.InitTagGroup("", []string{})
		relExpectedPath := directory + "/ec2_tagged.json"
		cfnBlocks := []structure.IBlock{
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
		f, _ := ioutil.TempFile(directory, "cfn.*.json")
		err := WriteJsonFile(readFilePath, cfnBlocks, f.Name(), nil)
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

func TestMapBracketsInFile(t *testing.T) {
	t.Run("one line", func(t *testing.T) {
		str := "{}[] not brackets ["
		expected := []types.Brackets{
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 1},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 2},
			{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 3},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 18},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("one line, nested", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		expected := []types.Brackets{
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 7},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 11},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 20},
			{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 26},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 27},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("multiple lines", func(t *testing.T) {
		str := "{\n}[] not \nbrackets \n["
		expected := []types.Brackets{
			{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 2, CharIndex: 2},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 3},
			{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 4},
			{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 4, CharIndex: 21},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
}

func TestGetBracketsPairs(t *testing.T) {
	t.Run("one line, no nesting", func(t *testing.T) {
		str := "{}[] not brackets"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)
		expectedPairs := map[int][]types.Brackets{
			0: {
				{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
				{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 1}},
			2: {
				{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 2},
				{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 3}},
		}

		assert.Equal(t, expectedPairs, actualPairs)
	})

	t.Run("one line, nesting", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)

		expectedPairs := map[int][]types.Brackets{
			0: {
				{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
				{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 27}},
			7: {
				{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 7},
				{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 11}},
			20: {
				{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 20},
				{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 1, CharIndex: 26}},
		}
		for index, pair := range expectedPairs {
			actualPair, ok := actualPairs[index]
			if !ok {
				t.Errorf("expected to get pair in index %d", index)
			}
			assert.Equal(t, pair, actualPair)
		}
	})
	t.Run("multiple lines with nesting", func(t *testing.T) {
		str := "{bana: {1:1},\n bana2:[1,2,3]}"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)

		expectedPairs := map[int][]types.Brackets{
			0: {
				{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 0},
				{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 2, CharIndex: 28}},
			7: {
				{Type: types.OpenBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 7},
				{Type: types.CloseBrackets, Shape: types.CurlyBrackets, Line: 1, CharIndex: 11}},
			21: {
				{Type: types.OpenBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 21},
				{Type: types.CloseBrackets, Shape: types.SquareBrackets, Line: 2, CharIndex: 27}},
		}
		for index, pair := range expectedPairs {
			actualPair, ok := actualPairs[index]
			if !ok {
				t.Errorf("expected to get pair in index %d", index)
			}
			assert.Equal(t, pair, actualPair)
		}
	})
}
