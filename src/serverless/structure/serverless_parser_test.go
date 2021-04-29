package structure

import (
	"bridgecrewio/yor/src/common"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
		directory := "../../../tests/serverless/resources"
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
		assert.Equal(t, common.Lines{Start: 13, End: 18}, func1Block.GetLines())
		assert.Equal(t, "myFunction", func1Block.GetResourceID())

		existingTag := func1Block.GetExistingTags()[0]
		assert.Equal(t, "TAG1_FUNC", existingTag.GetKey())
		assert.Equal(t, "Func1 Tag Value", existingTag.GetValue())
	})

}

func compareLines(t *testing.T, expected map[string]*common.Lines, actual map[string]*common.Lines) {
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
		directory := "../../../tests/serverless/resources"
		slsFilepath, _ := filepath.Abs(strings.Join([]string{directory, "serverless.yml"}, "/"))
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		slsBlocks, err := slsParser.ParseFile(slsFilepath)
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		assert.Equal(t, 2, len(slsBlocks))
		func1Block := slsBlocks[0]
		expected := map[string]*common.Lines{
			"myFunction": {Start: 13, End: 18},
		}
		func1Lines := func1Block.GetLines()
		compareLines(t, expected, map[string]*common.Lines{"myFunction": &func1Lines})
	})

	t.Run("test multiple resources", func(t *testing.T) {
		directory := "../../../tests/serverless/resources"
		slsFilepath, _ := filepath.Abs(strings.Join([]string{directory, "serverless.yml"}, "/"))
		slsParser := ServerlessParser{}
		slsParser.Init(directory, nil)
		slsBlocks, err := slsParser.ParseFile(slsFilepath)
		func1Block := slsBlocks[0]
		func1Lines := func1Block.GetLines()

		func2Block := slsBlocks[1]
		func2Lines := func2Block.GetLines()

		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		expected := map[string]*common.Lines{
			"myFunction":  {Start: 13, End: 18},
			"myFunction2": {Start: 19, End: 25},
		}
		compareLines(t, expected, map[string]*common.Lines{"myFunction": &func1Lines, "myFunction2": &func2Lines})
	})
}
