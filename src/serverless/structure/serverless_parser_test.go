package structure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/simple"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/stretchr/testify/assert"
)

func TestServerlessParser_ParseFile(t *testing.T) {
	t.Run("parse serverless file", func(t *testing.T) {
		path, err := os.Getwd()
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		err = os.Chdir(path)
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		directory := "../../../tests/serverless/resources/tags_exist"
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		slsFilepath, _ := filepath.Abs(filepath.Join(slsParser.YamlParser.RootDir, "serverless.yml"))
		expectedSlsFilepath, _ := filepath.Abs(filepath.Join(slsParser.YamlParser.RootDir, "serverless_expected.yaml"))
		slsBlocks, err := slsParser.ParseFile(slsFilepath)
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		assert.Equal(t, 2, len(slsBlocks))
		var func1Block *ServerlessBlock
		var func2Block *ServerlessBlock
		for _, block := range slsBlocks {
			castedBlock := block.(*ServerlessBlock)
			if castedBlock.Name == "myFunction" {
				func1Block = castedBlock
				assert.Equal(t, structure.Lines{Start: 13, End: 18}, func1Block.GetLines())
				assert.Equal(t, "myFunction", func1Block.GetResourceID())

				expectedTags := []tags.ITag{
					&tags.Tag{Key: "TAG1_FUNC", Value: "Func1 Tag Value"},
					&tags.Tag{Key: "TAG2_FUNC", Value: "Func1 Tag2 Value"},
				}

				assert.ElementsMatch(t, expectedTags, func1Block.GetExistingTags())
			} else if castedBlock.Name == "myFunction2" {
				func2Block = castedBlock
				assert.Equal(t, structure.Lines{Start: 19, End: 24}, func2Block.GetLines())
				assert.Equal(t, "myFunction2", func2Block.GetResourceID())

				expectedTags := []tags.ITag{
					&tags.Tag{Key: "TAG1_FUNC", Value: "Func2 Tag Value"},
					&tags.Tag{Key: "TAG2_FUNC", Value: "Func2 Tag2 Value"},
				}

				assert.ElementsMatch(t, expectedTags, func2Block.GetExistingTags())
			}
		}
		assert.NotNil(t, func1Block)
		assert.NotNil(t, func2Block)
		f, _ := os.CreateTemp(directory, "serverless.*.yaml")
		_ = slsParser.WriteFile(slsFilepath, slsBlocks, f.Name(), false)

		expected, _ := os.ReadFile(expectedSlsFilepath)
		actual, _ := os.ReadFile(f.Name())

		assert.Equal(t, string(expected), string(actual))

		defer func() { _ = os.Remove(f.Name()) }()
	})
}

func compareLines(t *testing.T, expected map[string]*structure.Lines, actual map[string]*structure.Lines) {
	for resourceName := range expected {
		actualLines := actual[resourceName]
		if actualLines == nil {
			t.Errorf("expected %s to be in resources mapping", resourceName)
		}
		expctedLines := expected[resourceName]
		assert.Equal(t, expctedLines, actualLines)
	}
}

func Test_mapResourcesLineYAML(t *testing.T) {
	t.Run("test single resource", func(t *testing.T) {
		directory := "../../../tests/serverless/resources/tags_exist"
		slsFilepath, _ := filepath.Abs(strings.Join([]string{directory, "serverless.yml"}, "/"))
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		slsBlocks, err := slsParser.ParseFile(slsFilepath)
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		var func1Block *ServerlessBlock
		assert.Equal(t, 2, len(slsBlocks))
		for _, block := range slsBlocks {
			castedBlock := block.(*ServerlessBlock)
			if castedBlock.Name == "myFunction" {
				func1Block = castedBlock
			}
		}
		expected := map[string]*structure.Lines{
			"myFunction": {Start: 13, End: 18},
		}
		func1Lines := func1Block.GetLines()
		compareLines(t, expected, map[string]*structure.Lines{"myFunction": &func1Lines})
	})

	t.Run("test multiple resources", func(t *testing.T) {
		directory := "../../../tests/serverless/resources/tags_exist"
		slsFilepath, _ := filepath.Abs(strings.Join([]string{directory, "serverless.yml"}, "/"))
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		var func1Block, func2Block *ServerlessBlock
		slsBlocks, err := slsParser.ParseFile(slsFilepath)
		for _, block := range slsBlocks {
			castedBlock := block.(*ServerlessBlock)
			if castedBlock.Name == "myFunction" {
				func1Block = castedBlock
			} else {
				func2Block = castedBlock
			}

		}
		func1Lines := func1Block.GetLines()
		func2Lines := func2Block.GetLines()

		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		expected := map[string]*structure.Lines{
			"myFunction":  {Start: 13, End: 18},
			"myFunction2": {Start: 19, End: 24},
		}
		compareLines(t, expected, map[string]*structure.Lines{"myFunction": &func1Lines, "myFunction2": &func2Lines})
	})

	t.Run("test multiple resources no tags", func(t *testing.T) {
		directory := "../../../tests/serverless/resources/no_tags"
		slsFilepath, _ := filepath.Abs(strings.Join([]string{directory, "serverless.yml"}, "/"))
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		var func1Block, func2Block *ServerlessBlock
		slsBlocks, err := slsParser.ParseFile(slsFilepath)
		for _, block := range slsBlocks {
			castedBlock := block.(*ServerlessBlock)
			if castedBlock.Name == "myFunction" {
				func1Block = castedBlock
			} else {
				func2Block = castedBlock
			}

		}
		func1Lines := func1Block.GetLines()
		func2Lines := func2Block.GetLines()

		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		expected := map[string]*structure.Lines{
			"myFunction":  {Start: 13, End: 15},
			"myFunction2": {Start: 16, End: 18},
		}
		compareLines(t, expected, map[string]*structure.Lines{"myFunction": &func1Lines, "myFunction2": &func2Lines})
	})

	t.Run("test try parse non serverless file name", func(t *testing.T) {
		directory := "../../../tests/serverless/resources/non_serverless"
		slsFilepath, _ := filepath.Abs(strings.Join([]string{directory, "file.yml"}, "/"))
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		parsedBlocks, _ := slsParser.ParseFile(slsFilepath)
		if parsedBlocks != nil {
			t.Fail()
		}
	})

	t.Run("test SLS writing", func(t *testing.T) {
		directory := "../../../tests/serverless/resources/no_tags"
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		readFilePath := directory + "/serverless.yml"
		tagGroup := simple.TagGroup{}
		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
		}
		tagGroup.SetTags(extraTags)
		tagGroup.InitTagGroup("", []string{}, []string{})
		writeFilePath := directory + "/serverless_tagged.yml"
		slsBlocks, err := slsParser.ParseFile(readFilePath)
		for _, block := range slsBlocks {
			err := tagGroup.CreateTagsForBlock(block)
			if err != nil {
				t.Fail()
			}
		}
		if err != nil {
			t.Fail()
		}
		f, _ := os.CreateTemp(directory, "serverless.*.yaml")
		err = slsParser.WriteFile(readFilePath, slsBlocks, f.Name(), false)
		if err != nil {
			t.Fail()
		}
		expectedFilePath, _ := filepath.Abs(writeFilePath)
		actualFilePath, _ := filepath.Abs(f.Name())
		expected, _ := os.ReadFile(expectedFilePath)
		actualOutput, _ := os.ReadFile(actualFilePath)
		assert.Equal(t, string(expected), string(actualOutput))
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Fail()
			}
		}(f.Name())

	})
}
