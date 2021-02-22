package structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloudformationParser_ParseFile(t *testing.T) {
	t.Run("parse ebs file", func(t *testing.T) {
		directory := "../../tests/cloudformation/resources"
		cfnParser := CloudformationParser{}
		cfnParser.Init(directory, nil)
		cfnBlocks, err := cfnParser.ParseFile(directory + "/ebs.yaml")
		if err != nil {
			t.Errorf("ParseFile() error = %v", err)
			return
		}
		assert.Equal(t, 1, len(cfnBlocks))
		newVolumeBlock := cfnBlocks[0]
		assert.Equal(t, []int{4, 14}, newVolumeBlock.GetLines())
		assert.Equal(t, "NewVolume", newVolumeBlock.GetResourceID())

		existingTag := newVolumeBlock.GetExistingTags()[0]
		assert.Equal(t, "MyTag", existingTag.GetKey())
		assert.Equal(t, "TagValue", existingTag.GetValue())
	})

}

func compareLines(t *testing.T, expected map[string][]int, actual map[string][]int) {
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
		filePath := "../../tests/cloudformation/resources/ebs.yaml"
		resourcesNames := []string{"NewVolume"}
		expected := map[string][]int{
			"NewVolume": {4, 14},
		}
		actual := mapResourcesLineYAML(filePath, resourcesNames)
		compareLines(t, expected, actual)
	})

	t.Run("test multiple resources", func(t *testing.T) {
		filePath := "../../tests/cloudformation/resources/ec2_untagged.yaml"
		resourcesNames := []string{"EC2InstanceResource0", "EC2InstanceResource1", "EC2LaunchTemplateResource0", "EC2LaunchTemplateResource1"}
		expected := map[string][]int{
			"EC2InstanceResource0":       {3, 6},
			"EC2InstanceResource1":       {7, 16},
			"EC2LaunchTemplateResource0": {17, 21},
			"EC2LaunchTemplateResource1": {22, 32},
		}
		actual := mapResourcesLineYAML(filePath, resourcesNames)
		compareLines(t, expected, actual)
	})
}
