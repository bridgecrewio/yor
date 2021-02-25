package structure

import (
	"bridgecrewio/yor/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloudformationParser_ParseFile(t *testing.T) {
	t.Run("parse ebs file", func(t *testing.T) {
		directory := "../../tests/cloudformation/resources/ebs"
		cfnParser := CloudformationParser{}
		cfnParser.Init(directory, nil)
		cfnBlocks, err := cfnParser.ParseFile(directory + "/ebs.yaml")
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		assert.Equal(t, 1, len(cfnBlocks))
		newVolumeBlock := cfnBlocks[0]
		assert.Equal(t, common.Lines{Start: 4, End: 14}, newVolumeBlock.GetLines())
		assert.Equal(t, "NewVolume", newVolumeBlock.GetResourceID())

		existingTag := newVolumeBlock.GetExistingTags()[0]
		assert.Equal(t, "MyTag", existingTag.GetKey())
		assert.Equal(t, "TagValue", existingTag.GetValue())
		assert.Equal(t, 4, cfnParser.fileToResourcesLines[directory+"/ebs.yaml"].Start)
		assert.Equal(t, 14, cfnParser.fileToResourcesLines[directory+"/ebs.yaml"].End)
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
		filePath := "../../tests/cloudformation/resources/ebs/ebs.yaml"
		resourcesNames := []string{"NewVolume"}
		expected := map[string]*common.Lines{
			"NewVolume": {Start: 4, End: 14},
		}
		actual := mapResourcesLineYAML(filePath, resourcesNames)
		compareLines(t, expected, actual)
	})

	t.Run("test multiple resources", func(t *testing.T) {
		filePath := "../../tests/cloudformation/resources/ec2_untagged.yaml"
		resourcesNames := []string{"EC2InstanceResource0", "EC2InstanceResource1", "EC2LaunchTemplateResource0", "EC2LaunchTemplateResource1"}
		expected := map[string]*common.Lines{
			"EC2InstanceResource0":       {3, 6},
			"EC2InstanceResource1":       {7, 16},
			"EC2LaunchTemplateResource0": {17, 21},
			"EC2LaunchTemplateResource1": {22, 32},
		}
		actual := mapResourcesLineYAML(filePath, resourcesNames)
		compareLines(t, expected, actual)
	})
}
