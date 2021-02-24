package gitservice

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

const TerragoatURL = "https://github.com/bridgecrewio/terragoat.git"

func TestNewGitService(t *testing.T) {
	t.Run("Get correct organization and repo name", func(t *testing.T) {
		terragoatPath := CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		gitService, err := NewGitService(terragoatPath)
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Get correct organization and repo name when in non-root dir", func(t *testing.T) {
		terragoatPath := CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)
		gitService, err := NewGitService(terragoatPath + "/aws")
		if err != nil {
			t.Errorf("could not initialize git service becauses %s", err)
		}
		assert.Equal(t, "bridgecrewio", gitService.GetOrganization())
		assert.Equal(t, "terragoat", gitService.GetRepoName())
	})

	t.Run("Fail if gotten to root dir", func(t *testing.T) {
		terragoatPath := CloneRepo(TerragoatURL)
		defer os.RemoveAll(terragoatPath)

		terragoatPath = filepath.Dir(filepath.Dir(terragoatPath))
		gitService, err := NewGitService(terragoatPath)
		assert.NotNil(t, err)
		assert.Nil(t, gitService)
	})
}

func CloneRepo(repoPath string) string {
	dir, err := ioutil.TempDir("", "temp-repo")
	if err != nil {
		log.Fatal(err)
	}

	// Clones the repository into the given dir, just as a normal git clone does
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoPath,
	})

	if err != nil {
		log.Fatal(err)
	}

	return dir
}
