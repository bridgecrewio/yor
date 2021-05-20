package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerrraformModule(t *testing.T) {
	t.Run("Test TF Module remote https logic", func(t *testing.T) {
		isRemote := isRemoteModule("https://github.com/terraform-aws-modules/terraform-aws-vpc.git")
		assert.True(t, isRemote)
	})

	t.Run("Test TF Module remote github logic", func(t *testing.T) {
		isRemote := isRemoteModule("github.com/terraform-aws-modules/terraform-aws-vpc.git")
		assert.True(t, isRemote)
	})

	t.Run("Test TF Module remote git logic", func(t *testing.T) {
		isRemote := isRemoteModule("git@github.com:terraform-aws-modules/terraform-aws-vpc.git")
		assert.True(t, isRemote)
	})

	t.Run("Test TF Module local logic", func(t *testing.T) {
		localPath := "test/local/path"
		isRemote := isRemoteModule(localPath)
		assert.False(t, isRemote)
		isRegistry := isTerraformRegistryModule(localPath)
		assert.False(t, isRegistry)
	})

	t.Run("Test TF Registry Module logic", func(t *testing.T) {
		isRegistry := isTerraformRegistryModule("terraform-aws-modules/security-group/aws")
		assert.True(t, isRegistry)
	})
}
