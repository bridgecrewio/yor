package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/tagging/tags"
	"bridgecrewio/yor/tests/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagCreation(t *testing.T) {
	blame := utils.SetupBlame(t)
	t.Run("BcTraceTagCreation", func(t *testing.T) {
		tag := tags.YorTraceTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "yor_trace", tag.Key)
		assert.Equal(t, 36, len(tag.Value))
	})

	t.Run("GitOrgTagCreation", func(t *testing.T) {
		tag := tags.GitOrgTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_org", tag.Key)
		assert.Equal(t, utils.Org, tag.Value)
	})

	t.Run("GitRepoTagCreation", func(t *testing.T) {
		tag := tags.GitRepoTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_repo", tag.Key)
		assert.Equal(t, utils.Repository, tag.Value)
	})

	t.Run("GitFileTagCreation", func(t *testing.T) {
		tag := tags.GitFileTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_file", tag.Key)
		assert.Equal(t, utils.FilePath, tag.Value)
	})

	t.Run("GitCommitTagCreation", func(t *testing.T) {
		tag := tags.GitCommitTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_commit", tag.Key)
		assert.Equal(t, utils.CommitHash1, tag.Value)
	})

	t.Run("GitLastModifiedAtCreation", func(t *testing.T) {
		tag := tags.GitLastModifiedAtTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_last_modified_at", tag.Key)
		assert.Equal(t, "2020-03-28 21:42:46 +0000 UTC", tag.Value)
	})

	t.Run("GitLastModifiedByCreation", func(t *testing.T) {
		tag := tags.GitLastModifiedByTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_last_modified_by", tag.Key)
		assert.Equal(t, "schosterbarak@gmail.com", tag.Value)
	})

	t.Run("GitModifiersCreation", func(t *testing.T) {
		tag := tags.GitModifiersTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_modifiers", tag.Key)
		assert.Equal(t, "jonjozwiak/schosterbarak", tag.Value)
	})

}

func EvaluateTag(t *testing.T, tag tags.ITag, blame git_service.GitBlame) {
	tag.Init()
	err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
}
