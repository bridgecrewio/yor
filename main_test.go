package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadExternalTags(t *testing.T) {
	t.Run("load local plugins", func(t *testing.T) {
		pluginDir := "yor_plugins/example"
		fmt.Printf("please make sure you have .so file in %s. if not, run the following command: \n", pluginDir)
		fmt.Printf("go build -gcflags=\"all=-N -l\" -buildmode=plugin -o %s/extra_tags.so %s/*.go\n", pluginDir, pluginDir)
		gotTags, err := loadExternalTags(pluginDir)
		if err != nil {
			t.Errorf("loadExternalTags() error = %v", err)
			return
		}
		expectedTags := map[string]string{"yor_checkov": "checkov", "yor_terragoat": "terragoat"}
		assert.Equal(t, len(expectedTags), len(gotTags))

		for _, tag := range gotTags {
			key := tag.GetKey()
			value := tag.GetValue()
			assert.Equal(t, expectedTags[key], value)
		}
	})
}
