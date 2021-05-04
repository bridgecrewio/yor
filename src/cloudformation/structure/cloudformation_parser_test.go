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
	"github.com/bridgecrewio/yor/src/common/utils"
	"github.com/stretchr/testify/assert"
)

func TestCloudformationParser_ParseFile(t *testing.T) {
	t.Run("parse ebs file", func(t *testing.T) {
		directory := "../../../tests/cloudformation/resources/ebs"
		cfnParser := CloudformationParser{}
		cfnParser.Init(directory, nil)
		cfnBlocks, err := cfnParser.ParseFile(directory + "/ebs.yaml")
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		assert.Equal(t, 1, len(cfnBlocks))
		newVolumeBlock := cfnBlocks[0]
		assert.Equal(t, structure.Lines{Start: 4, End: 14}, newVolumeBlock.GetLines())
		assert.Equal(t, "NewVolume", newVolumeBlock.GetResourceID())

		existingTag := newVolumeBlock.GetExistingTags()[0]
		assert.Equal(t, "MyTag", existingTag.GetKey())
		assert.Equal(t, "TagValue", existingTag.GetValue())
		assert.Equal(t, 4, cfnParser.FileToResourcesLines[directory+"/ebs.yaml"].Start)
		assert.Equal(t, 14, cfnParser.FileToResourcesLines[directory+"/ebs.yaml"].End)
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
		filePath := "../../../tests/cloudformation/resources/ebs/ebs.yaml"
		resourcesNames := []string{"NewVolume"}
		expected := map[string]*structure.Lines{
			"NewVolume": {Start: 4, End: 14},
		}
		actual := MapResourcesLineYAML(filePath, resourcesNames)
		compareLines(t, expected, actual)
	})

	t.Run("test multiple resources", func(t *testing.T) {
		filePath := "../../../tests/cloudformation/resources/ec2_untagged/ec2_untagged.yaml"
		resourcesNames := []string{"EC2InstanceResource0", "EC2InstanceResource1", "EC2LaunchTemplateResource0", "EC2LaunchTemplateResource1"}
		expected := map[string]*structure.Lines{
			"EC2InstanceResource0":       {Start: 3, End: 6},
			"EC2InstanceResource1":       {Start: 7, End: 16},
			"EC2LaunchTemplateResource0": {Start: 17, End: 21},
			"EC2LaunchTemplateResource1": {Start: 22, End: 32},
		}
		actual := MapResourcesLineYAML(filePath, resourcesNames)
		compareLines(t, expected, actual)
	})

	t.Run("test CFN writing", func(t *testing.T) {
		directory := "../../../tests/cloudformation/resources/ebs"
		f, _ := ioutil.TempFile(directory, "temp.*.yaml")
		cfnParser := CloudformationParser{}
		cfnParser.Init(directory, nil)
		readFilePath := directory + "/ebs.yaml"
		tagGroup := simple.TagGroup{}
		extraTags := []tags.ITag{
			&tags.Tag{
				Key:   "new_tag",
				Value: "new_value",
			},
		}
		tagGroup.SetTags(extraTags)
		tagGroup.InitTagGroup("", []string{})
		writeFilePath := directory + "/ebs_tagged.yaml"
		cfnBlocks, err := cfnParser.ParseFile(readFilePath)
		for _, block := range cfnBlocks {
			utils.CreateTagsForBlock(&tagGroup, block)
		}
		if err != nil {
			t.Fail()
		}
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			t.Fail()
		}
		err = cfnParser.WriteFile(readFilePath, cfnBlocks, f.Name())
		if err != nil {
			t.Fail()
		}
		var expectedHandler, actualHandler *os.File
		expectedAbs, _ := filepath.Abs(writeFilePath)
		actualAbs, _ := filepath.Abs(f.Name())
		expectedHandler, _ = os.OpenFile(expectedAbs, os.O_RDWR, 0755)
		actualHandler, _ = os.OpenFile(actualAbs, os.O_RDWR|os.O_CREATE, 0755)
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
