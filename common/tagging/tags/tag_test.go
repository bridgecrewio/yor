package tags

import (
	"bridgecrewio/yor/common/gitservice"
	"bridgecrewio/yor/tests/utils/blameutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagCreation(t *testing.T) {
	blame := blameutils.SetupBlame(t)
	t.Run("BcTraceTagCreation", func(t *testing.T) {
		tag := YorTraceTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "yor_trace", tag.Key)
		assert.Equal(t, 36, len(tag.Value))
	})

	t.Run("GitOrgTagCreation", func(t *testing.T) {
		tag := GitOrgTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_org", tag.Key)
		assert.Equal(t, blameutils.Org, tag.Value)
	})

	t.Run("GitRepoTagCreation", func(t *testing.T) {
		tag := GitRepoTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_repo", tag.Key)
		assert.Equal(t, blameutils.Repository, tag.Value)
	})

	t.Run("GitFileTagCreation", func(t *testing.T) {
		tag := GitFileTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_file", tag.Key)
		assert.Equal(t, blameutils.FilePath, tag.Value)
	})

	t.Run("GitCommitTagCreation", func(t *testing.T) {
		tag := GitCommitTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_commit", tag.Key)
		assert.Equal(t, blameutils.CommitHash1, tag.Value)
	})

	t.Run("GitLastModifiedAtCreation", func(t *testing.T) {
		tag := GitLastModifiedAtTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_last_modified_at", tag.Key)
		assert.Equal(t, "2020-03-28 21:42:46", tag.Value)
	})

	t.Run("GitLastModifiedByCreation", func(t *testing.T) {
		tag := GitLastModifiedByTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_last_modified_by", tag.Key)
		assert.Equal(t, "schosterbarak@gmail.com", tag.Value)
	})

	t.Run("GitModifiersCreation", func(t *testing.T) {
		tag := GitModifiersTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_modifiers", tag.Key)
		assert.Equal(t, "jonjozwiak/schosterbarak", tag.Value)
	})

}

func EvaluateTag(t *testing.T, tag ITag, blame gitservice.GitBlame) {
	tag.Init()
	err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
}
