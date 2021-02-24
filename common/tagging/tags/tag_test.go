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
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "yor_trace", valueTag.GetKey())
		assert.Equal(t, 36, len(valueTag.GetValue()))
	})

	t.Run("GitOrgTagCreation", func(t *testing.T) {
		tag := GitOrgTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_org", valueTag.GetKey())
		assert.Equal(t, blameutils.Org, valueTag.GetValue())
	})

	t.Run("GitRepoTagCreation", func(t *testing.T) {
		tag := GitRepoTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_repo", valueTag.GetKey())
		assert.Equal(t, blameutils.Repository, valueTag.GetValue())
	})

	t.Run("GitFileTagCreation", func(t *testing.T) {
		tag := GitFileTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_file", valueTag.GetKey())
		assert.Equal(t, blameutils.FilePath, valueTag.GetValue())
	})

	t.Run("GitCommitTagCreation", func(t *testing.T) {
		tag := GitCommitTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_commit", valueTag.GetKey())
		assert.Equal(t, blameutils.CommitHash1, valueTag.GetValue())
	})

	t.Run("GitLastModifiedAtCreation", func(t *testing.T) {
		tag := GitLastModifiedAtTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_last_modified_at", valueTag.GetKey())
		assert.Equal(t, "2020-03-28 21:42:46", valueTag.GetValue())
	})

	t.Run("GitLastModifiedByCreation", func(t *testing.T) {
		tag := GitLastModifiedByTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_last_modified_by", valueTag.GetKey())
		assert.Equal(t, "schosterbarak@gmail.com", valueTag.GetValue())
	})

	t.Run("GitModifiersCreation", func(t *testing.T) {
		tag := GitModifiersTag{}
		valueTag := EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_modifiers", valueTag.GetKey())
		assert.Equal(t, "jonjozwiak/schosterbarak", valueTag.GetValue())
	})

}

func EvaluateTag(t *testing.T, tag ITag, blame gitservice.GitBlame) ITag {
	tag.Init()
	newTag, err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &Tag{}, newTag)

	return newTag
}
