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
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		defer os.RemoveAll(terragoatPath)

		gitService, err := NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Get correct organization and repo name when in non-root dir", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		defer os.RemoveAll(terragoatPath)
		gitService, err := NewGitService(terragoatPath + "/aws")
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Get correct organization and repo name from deeper gitlab", func(t *testing.T) {
		gitlabPath := utils.CloneRepo("https://gitlab.com/gitlab-org/configure/examples/gitlab-terraform-aws.git", "4e45d0983ec157376b3389f08e565acdc6f49eee")
		defer os.RemoveAll(gitlabPath)
		gitService, err := NewGitService(gitlabPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "gitlab-org", gitService.GetOrganization())
		assert.Equal(t, "configure/examples/gitlab-terraform-aws", gitService.GetRepoName())
	})

	t.Run("Fail if gotten to root dir", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		defer os.RemoveAll(terragoatPath)

		terragoatPath = filepath.Dir(filepath.Dir(terragoatPath))
		gitService, err := NewGitService(terragoatPath)
		assert.NotNil(t, err)
		assert.Nil(t, gitService)
	})

	t.Run("Fail if gotten to root dir 2", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		defer os.RemoveAll(terragoatPath)

		terragoatPath = filepath.Dir(filepath.Dir(terragoatPath))
		gitService, err := NewGitService(terragoatPath)
		assert.NotNil(t, err)
		assert.Nil(t, gitService)
	})

	t.Run("Get blame for lines test", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		defer func() {
			_ = os.RemoveAll(terragoatPath)
		}()
		gitService, _ := NewGitService(terragoatPath)

		blame, _ := gitService.GetBlameForFileLines("terraform/aws/s3.tf", structure.Lines{Start: 1, End: 13})
		commit := blame.GetLatestCommit()
		assert.Equal(t, 13, len(blame.BlamesByLine))
		assert.Equal(t, "8ee70df3dcba70ab1bafafad5f7efe2a41cd53a3", commit.Hash.String())
	})

	t.Run("Get blame for lines test", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "52d676e1cd75a13f990f1caca51d6ad8c78858da")
		defer func() {
			_ = os.RemoveAll(terragoatPath)
		}()
		gitService, _ := NewGitService(terragoatPath)

		blame, _ := gitService.GetBlameForFileLines("terraform/aws/s3.tf", structure.Lines{Start: 1, End: 13})
		commit := blame.GetLatestCommit()
		assert.Equal(t, 13, len(blame.BlamesByLine))
		assert.Equal(t, "5c6b5d60a8aa63a5d37e60f15185d13a967f0542", commit.Hash.String())
		assert.Equal(t, "nimrodkor@users.noreply.github.com", commit.Author)
	})

	t.Run("Get correct relative file path", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		defer func() {
			_ = os.RemoveAll(terragoatPath)
		}()

		gitService, err := NewGitService(filepath.Join(terragoatPath, "terraform", "aws"))
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		targetPath := gitService.ComputeRelativeFilePath("aws/db-app.tf")
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
		assert.Equal(t, "terraform/aws/db-app.tf", targetPath)
	})

	t.Run("Get correct organization and repo name inside dir relative", func(t *testing.T) {
		terragoatPath := utils.CloneRepo(utils.TerragoatURL, "063dc2db3bb036160ed39d3705508ee8293a27c8")
		defer func() {
			_ = os.RemoveAll(terragoatPath)
		}()
		cwd, _ := os.Getwd()
		terragoatAbsPath := filepath.Join(terragoatPath, "terraform", "aws")
		relPath, _ := filepath.Rel(cwd, terragoatAbsPath)
		gitService, err := NewGitService(relPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		targetPath := gitService.ComputeRelativeFilePath("aws/db-app.tf")
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
		assert.Equal(t, "terraform/aws/db-app.tf", targetPath)
	})
}
