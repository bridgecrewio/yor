package gitservice

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/tests/utils"

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

	t.Run("Fail if gotten to root dir 2", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		terragoatPath = filepath.Dir(filepath.Dir(terragoatPath))
		gitService, err := NewGitService(terragoatPath)
		assert.NotNil(t, err)
		assert.Nil(t, gitService)
	})

	t.Run("Get blame for lines test", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL)
		defer func() {
			_ = os.RemoveAll(terragoatPath)
		}()
		gitService, _ := NewGitService(terragoatPath)

		blame, _ := gitService.GetBlameForFileLines("terraform/aws/s3.tf", structure.Lines{Start: 1, End: 13})
		commit := blame.GetLatestCommit()
		assert.Equal(t, 13, len(blame.BlamesByLine))
		assert.Equal(t, "5c6b5d60a8aa63a5d37e60f15185d13a967f0542", commit.Hash.String())
	})
}
