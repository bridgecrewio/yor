package gitservice

import (
	"bridgecrewio/yor/tests/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitService(t *testing.T) {
	t.Run("Get correct organization and repo name", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		gitService, err := NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Get correct organization and repo name when in non-root dir", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL)
		defer os.RemoveAll(terragoatPath)
		gitService, err := NewGitService(terragoatPath + "/aws")
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Fail if gotten to root dir", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		terragoatPath = filepath.Dir(filepath.Dir(terragoatPath))
		gitService, err := NewGitService(terragoatPath)
		assert.NotNil(t, err)
		assert.Nil(t, gitService)
	})
}
