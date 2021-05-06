package structure

import (
	"bufio"
	"io"
	"io/ioutil"
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
		slsFilepath, _ := filepath.Abs(strings.Join([]string{slsParser.YamlParser.RootDir, "serverless.yml"}, "/"))
		slsBlocks, err := slsParser.ParseFile(slsFilepath)
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		assert.Equal(t, 2, len(slsBlocks))
		var func1Block *ServerlessBlock
		for _, block := range slsBlocks {
			castedBlock := block.(*ServerlessBlock)
			if castedBlock.Name == "myFunction" {
				func1Block = castedBlock
			}
		}
		if func1Block == nil {
			assert.Fail(t, "Didn't find the function block")
		} else {
			assert.Equal(t, structure.Lines{Start: 14, End: 19}, func1Block.GetLines())
			assert.Equal(t, "myFunction", func1Block.GetResourceID())

			expectedTags := []tags.ITag{
				&tags.Tag{Key: "TAG1_FUNC", Value: "Func1 Tag Value"},
				&tags.Tag{Key: "TAG2_FUNC", Value: "Func1 Tag2 Value"},
			}

			assert.ElementsMatch(t, func1Block.GetExistingTags(), expectedTags)
		}
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
		for _, block := range slsBlocks {
			castedBlock := block.(*ServerlessBlock)
			if castedBlock.Name == "myFunction" {
				func1Block = castedBlock
			}
			assert.Equal(t, 2, len(slsBlocks))
			expected := map[string]*structure.Lines{
				"myFunction": {Start: 14, End: 19},
			}
			func1Lines := func1Block.GetLines()
			compareLines(t, expected, map[string]*structure.Lines{"myFunction": &func1Lines})
		}
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
			"myFunction":  {Start: 14, End: 19},
			"myFunction2": {Start: 20, End: 27},
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
			"myFunction":  {Start: 14, End: 16},
			"myFunction2": {Start: 17, End: 21},
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
		f, _ := ioutil.TempFile(directory, "serverless.*.yaml")
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
		tagGroup.InitTagGroup("", []string{})
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
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			t.Fail()
		}
		err = slsParser.WriteFile(readFilePath, slsBlocks, f.Name())
		if err != nil {
			t.Fail()
		}
		var expectedHandler, actualHandler *os.File
		expectedAbs, _ := filepath.Abs(writeFilePath)
		actualAbs, _ := filepath.Abs(f.Name())
		expectedHandler, _ = os.OpenFile(expectedAbs, os.O_RDWR, 0755)
		actualHandler, _ = os.OpenFile(actualAbs, os.O_RDWR|os.O_CREATE, 0755)
		defer expectedHandler.Close()
		defer actualHandler.Close()
		_, err = expectedHandler.Seek(0, io.SeekStart)
		if err != nil {
			t.Fail()
		}
		_, err = actualHandler.Seek(0, io.SeekStart)
		if err != nil {
			t.Fail()
		}
		defer expectedHandler.Close()
		defer actualHandler.Close()
		actualReader := bufio.NewScanner(actualHandler)
		expectedReader := bufio.NewScanner(expectedHandler)
		for actualReader.Scan() && expectedReader.Scan() {
			actualLine := actualReader.Text()
			expectedLine := expectedReader.Text()
			assert.Equal(t, strings.Trim(actualLine, " \n\t"), strings.Trim(expectedLine, " \n\t"))
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Fail()
			}
		}(f.Name())

	})
}
