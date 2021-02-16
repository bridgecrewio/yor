package runner

import (
	"bridgecrewio/yor/common/gitservice"
	"fmt"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/stretchr/testify/assert"
)

func Test_loadExternalTags(t *testing.T) {
	t.Run("load local plugins", func(t *testing.T) {
		pluginDir := "../../tests/yor_plugins/example"
		fmt.Printf("please make sure you have .so file in %s. if not, run the following command: \n", pluginDir)
		fmt.Printf("go build -gcflags=\"all=-N -l\" -buildmode=plugin -o %s/extra_tags.so %s/*.go\n", pluginDir, pluginDir)
		gotTags, err := loadExternalTags(pluginDir)
		if err != nil {
			t.Errorf("loadExternalTags() error = %v", err)
			return
		}
		expectedTags := map[string]string{"yor_checkov": "checkov", "git_owner": "bana"}
		assert.Equal(t, len(expectedTags), len(gotTags))
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		gitBlame := gitservice.GitBlame{
			GitOrg:        "bridgecrewio",
			GitRepository: "yor",
			BlamesByLine: map[int]*git.Line{0: {
				Author: "bana",
				Date:   now,
				Hash:   plumbing.NewHash("0"),
			}, 1: {Author: "shati",
				Date: yesterday,
				Hash: plumbing.NewHash("1")}}}
		for _, tag := range gotTags {
			tag.Init()
			err := tag.CalculateValue(&gitBlame)
			print(err)
			key := tag.GetKey()
			value := tag.GetValue()
			assert.Equal(t, expectedTags[key], value)
		}
	})
}
