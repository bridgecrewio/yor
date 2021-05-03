package structure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"

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
		func1Block := slsBlocks[0]
		assert.Equal(t, structure.Lines{Start: 14, End: 19}, func1Block.GetLines())
		assert.Equal(t, "myFunction", func1Block.GetResourceID())

		existingTag := func1Block.GetExistingTags()[0]
		assert.Equal(t, "TAG1_FUNC", existingTag.GetKey())
		assert.Equal(t, "Func1 Tag Value", existingTag.GetValue())
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
}
