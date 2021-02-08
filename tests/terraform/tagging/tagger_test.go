package tagging

import (
	"bridgecrewio/yor/terraform/tagging"
	"bridgecrewio/yor/tests/utils"
	"testing"
)

func TestTerraformTagger(t *testing.T) {
	blame := utils.SetupBlame(t)

	t.Run("test terraform tagger IsBlockTaggable", func(t *testing.T) {
		tagger := tagging.TerraformTagger{}

		_ = tagger
	})

	_ = blame
}
