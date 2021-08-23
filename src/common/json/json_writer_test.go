package json

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/awslabs/goformation/v5/cloudformation/s3"
	s3tags "github.com/awslabs/goformation/v5/cloudformation/tags"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"

	"github.com/stretchr/testify/assert"
)

func TestMapBracketsInFile(t *testing.T) {
	t.Run("one line", func(t *testing.T) {
		str := "{}[] not brackets ["
		expected := []Brackets{
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 1},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 2},
			{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 3},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 18},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("one line, nested", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		expected := []Brackets{
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 20},
			{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 26},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 27},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
	t.Run("multiple lines", func(t *testing.T) {
		str := "{\n}[] not \nbrackets \n["
		expected := []Brackets{
			{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
			{Type: CloseBrackets, Shape: CurlyBrackets, Line: 2, CharIndex: 2},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 3},
			{Type: CloseBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 4},
			{Type: OpenBrackets, Shape: SquareBrackets, Line: 4, CharIndex: 21},
		}
		actual := MapBracketsInString(str)
		assert.Equal(t, expected, actual)
	})
}

func TestTagWriting(t *testing.T) {
	t.Run("Test UpdateExistingTags", func(t *testing.T) {
		tagLinesList := []string{
			"[",
			"  {",
			`    "Value":"Reverse",`,
			`    "Key": "RK"`,
			"  },",
			"  {",
			`    "Key" : "DK",`,
			`    "Value" : "Direct"`,
			"  }",
			"]",
		}
		UpdateExistingTags(tagLinesList, []*tags.TagDiff{
			{Key: "RK", NewValue: "ReverseCorrect", PrevValue: "Reverse"},
			{Key: "DK", NewValue: "DirectCorrect", PrevValue: "Direct"},
		})

		assert.Equal(t, `    "Value": "ReverseCorrect",`, tagLinesList[2])
		assert.Equal(t, `    "Value": "DirectCorrect"`, tagLinesList[7])
	})
}

func TestGetBracketsPairs(t *testing.T) {
	t.Run("one line, no nesting", func(t *testing.T) {
		str := "{}[] not brackets"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)
		expectedPairs := map[int]BracketPair{
			0: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 1}},
			2: {
				Open:  Brackets{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 2},
				Close: Brackets{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 3}},
		}

		assert.Equal(t, expectedPairs, actualPairs)
	})

	t.Run("one line, nesting", func(t *testing.T) {
		str := "{bana: {1:1}, bana2:[1,2,3]}"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)

		expectedPairs := map[int]BracketPair{
			0: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 27}},
			7: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11}},
			20: {
				Open:  Brackets{Type: OpenBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 20},
				Close: Brackets{Type: CloseBrackets, Shape: SquareBrackets, Line: 1, CharIndex: 26}},
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

		expectedPairs := map[int]BracketPair{
			0: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 2, CharIndex: 28}},
			7: {
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11}},
			21: {
				Open:  Brackets{Type: OpenBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 21},
				Close: Brackets{Type: CloseBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 27}},
		}
		for index, pair := range expectedPairs {
			actualPair, ok := actualPairs[index]
			if !ok {
				t.Errorf("expected to get pair in index %d", index)
			}
			assert.Equal(t, pair, actualPair)
		}
	})
	t.Run("un even number of brackets", func(t *testing.T) {
		str := "}[]}"
		bracketsInFile := MapBracketsInString(str)
		actualPairs := GetBracketsPairs(bracketsInFile)

		assert.Equal(t, make(map[int]BracketPair), actualPairs)
	})
}

func Test_findIndent(t *testing.T) {
	type args struct {
		str        string
		charToStop byte
		startIndex int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no indent",
			args: args{
				str:        "parent:{child:[]}",
				charToStop: '{',
				startIndex: 0,
			},
			want: "",
		},
		{
			name: "single space indent",
			args: args{
				str:        "parent: {child:[]}",
				charToStop: '{',
				startIndex: 0,
			},
			want: " ",
		},
		{
			name: "single space indent with new line",
			args: args{
				str:        "parent: \n {child:[]}",
				charToStop: '{',
				startIndex: 0,
			},
			want: " ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findIndent(tt.args.str, tt.args.charToStop, tt.args.startIndex); got != tt.want {
				t.Errorf("findIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getJSONStr(t *testing.T) {
	type args struct {
		jsonData interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "map",
			args: args{
				jsonData: map[string]string{"parent": "child"},
			},
			want: "{\"parent\":\"child\"}",
		},
		{
			name: "can't marshal",
			args: args{
				jsonData: make(chan int),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getJSONStr(tt.args.jsonData); got != tt.want {
				t.Errorf("getJSONStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapResourcesLineJSON(t *testing.T) {
	type args struct {
		filePath      string
		resourceNames []string
	}
	tests := []struct {
		name      string
		args      args
		wantLines map[string]*structure.Lines
	}{
		{
			name: "ebs",
			args: args{
				filePath:      "../../../tests/cloudformation/resources/ebs/ebs.json",
				resourceNames: []string{"NewVolume"},
			},
			wantLines: map[string]*structure.Lines{"NewVolume": {
				Start: 5,
				End:   25,
			}},
		},
		{
			name: "file doesn't exists",
			args: args{
				filePath:      "../../../tests/cloudformation/resources/ebs/no_such_file.json",
				resourceNames: []string{"NewVolume"},
			},
			wantLines: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := MapResourcesLineJSON(tt.args.filePath, tt.args.resourceNames)
			if !reflect.DeepEqual(got, tt.wantLines) {
				t.Errorf("MapResourcesLineJSON() got = %v, want %v", got, tt.wantLines)
			}
		})
	}
}

func TestFindParentIdentifier(t *testing.T) {
	type args struct {
		str             string
		childIdentifier string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "find parent",
			args: args{
				str:             "{\"parent\":{\"child\": 3}}",
				childIdentifier: "child",
			},
			want: "parent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindParentIdentifier(tt.args.str, tt.args.childIdentifier); got != tt.want {
				t.Errorf("FindParentIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findJSONKeyIndex(t *testing.T) {
	type args struct {
		str string
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "key exist",
			args: args{
				str: "{\"parent\":{\"child\": 3}}",
				key: "parent",
			},
			want: 1,
		},
		{
			name: "key doesn't exist",
			args: args{
				str: "{\"parent\":{\"child\": 3}}",
				key: "parent2",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findJSONKeyIndex(tt.args.str, tt.args.key); got != tt.want {
				t.Errorf("findJSONKeyIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getKeyIndex(t *testing.T) {
	type args struct {
		str        string
		key        string
		linesRange *structure.Lines
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "key index with multiple lines",
			args: args{
				str: "{\n\"parent\":{\"child\": 3}}",
				key: "parent",
				linesRange: &structure.Lines{
					Start: 1,
					End:   1,
				},
			},
			want: 2,
		},
		{
			name: "key index with multiple lines 1",
			args: args{
				str: "{\n\"parent\":{\n\"child\": 3}}",
				key: "parent",
				linesRange: &structure.Lines{
					Start: 1,
					End:   2,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getKeyIndex(tt.args.str, tt.args.key, tt.args.linesRange); got != tt.want {
				t.Errorf("getKeyIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindWrappingBrackets(t *testing.T) {
	type args struct {
		allBracketPairs  map[int]BracketPair
		innerBracketPair BracketPair
	}
	tests := []struct {
		name string
		args args
		want BracketPair
	}{
		{
			name: "",
			args: args{
				allBracketPairs: map[int]BracketPair{
					0: {
						Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
						Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 2, CharIndex: 28}},
					7: {
						Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
						Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11}},
					21: {
						Open:  Brackets{Type: OpenBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 21},
						Close: Brackets{Type: CloseBrackets, Shape: SquareBrackets, Line: 2, CharIndex: 27}},
				},
				innerBracketPair: BracketPair{
					Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 7},
					Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 11}},
			},
			want: BracketPair{
				Open:  Brackets{Type: OpenBrackets, Shape: CurlyBrackets, Line: 1, CharIndex: 0},
				Close: Brackets{Type: CloseBrackets, Shape: CurlyBrackets, Line: 2, CharIndex: 28}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindWrappingBrackets(tt.args.allBracketPairs, tt.args.innerBracketPair); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindWrappingBrackets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func writeJSONTestHelper(t *testing.T, directory string, testFileName string, existingTags []tags.Tag) {
	readFilePath := directory + "/" + testFileName + ".json"
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
	iTags := make([]tags.ITag, 0)
	for _, tag := range existingTags {
		tag := &tags.Tag{Key: tag.Key, Value: tag.Value}
		iTags = append(iTags, tag)
	}

	blocks := []structure.IBlock{
		&structure.Block{
			FilePath:    readFilePath,
			ExitingTags: iTags,
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
	f, _ := ioutil.TempFile(directory, testFileName+".*.json")
	_, fileBracketsMapping := MapResourcesLineJSON(readFilePath, []string{"S3Bucket"})
	err := WriteJSONFile(readFilePath, blocks, f.Name(), fileBracketsMapping)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	expectedFilePath := filepath.Join(directory, testFileName+"_expected.json")
	actualFilePath, _ := filepath.Abs(f.Name())
	expected, _ := ioutil.ReadFile(expectedFilePath)
	actualOutput, _ := ioutil.ReadFile(actualFilePath)
	assert.Equal(t, string(expected), string(actualOutput))
	defer func() {
		_ = os.Remove(f.Name())
	}()
}

func TestCFNWriting(t *testing.T) {
	t.Run("test CFN writing untagged", func(t *testing.T) {
		directory := "../../../tests/cloudformation/resources/no_tags"
		writeJSONTestHelper(t, directory, "base", []tags.Tag{})
	})
	t.Run("test CFN writing tagged", func(t *testing.T) {
		directory := "../../../tests/cloudformation/resources/json"
		writeJSONTestHelper(t, directory, "base", []tags.Tag{{Key: "old_tag", Value: "old_value"}})
	})
}
