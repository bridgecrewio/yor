package git_service

import (
	"bridgecrewio/yor/git_service"
	"github.com/stretchr/testify/assert"
	"testing"
)

const TerragoatUrl = "https://github.com/bridgecrewio/terragoat.git"

func TestGetBlameForFileLines(t *testing.T) {
	t.Run("mock test", func(t *testing.T) {
		gitService, err := git_service.NewGitService("")
		if err != nil {
			t.Errorf("could not initialize repository becauses %s", err)
		}
		blame, err := gitService.GetBlameForFileLines("", nil, "")
		assert.Nil(t, blame)
		assert.Nil(t, err)

	})
}
