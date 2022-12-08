package gittag

import (
	"os"
	"testing"

	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/tests/utils/blameutils"

	"github.com/stretchr/testify/assert"
)

func TestTagCreation(t *testing.T) {
	blame := blameutils.SetupBlame(t)
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

	t.Run("Tag description tests", func(t *testing.T) {
		tag := tags.Tag{}
		defaultDescription := tag.GetDescription()
		cwd, _ := os.Getwd()
		g := TagGroup{}
		g.InitTagGroup(cwd, nil, nil)
		for _, tag := range g.GetTags() {
			assert.NotEqual(t, defaultDescription, tag.GetDescription())
			assert.NotEqual(t, "", tag.GetDescription())
		}
	})

}

func TestTagCreationWithPrefix(t *testing.T) {
	blame := blameutils.SetupBlame(t)
	t.Run("GitOrgTagCreation", func(t *testing.T) {
		tag := GitOrgTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_git_org", valueTag.GetKey())
		assert.Equal(t, blameutils.Org, valueTag.GetValue())
	})

	t.Run("GitRepoTagCreation", func(t *testing.T) {
		tag := GitRepoTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_git_repo", valueTag.GetKey())
		assert.Equal(t, blameutils.Repository, valueTag.GetValue())
	})

	t.Run("GitFileTagCreation", func(t *testing.T) {
		tag := GitFileTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_git_file", valueTag.GetKey())
		assert.Equal(t, blameutils.FilePath, valueTag.GetValue())
	})

	t.Run("GitCommitTagCreation", func(t *testing.T) {
		tag := GitCommitTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_git_commit", valueTag.GetKey())
		assert.Equal(t, blameutils.CommitHash1, valueTag.GetValue())
	})

	t.Run("GitLastModifiedAtCreation", func(t *testing.T) {
		tag := GitLastModifiedAtTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_git_last_modified_at", valueTag.GetKey())
		assert.Equal(t, "2020-03-28 21:42:46", valueTag.GetValue())
	})

	t.Run("GitLastModifiedByCreation", func(t *testing.T) {
		tag := GitLastModifiedByTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_git_last_modified_by", valueTag.GetKey())
		assert.Equal(t, "schosterbarak@gmail.com", valueTag.GetValue())
	})

	t.Run("GitModifiersCreation", func(t *testing.T) {
		tag := GitModifiersTag{}
		valueTag := EvaluateTagWithPrefix(t, &tag, blame, "prefix_")
		assert.Equal(t, "prefix_git_modifiers", valueTag.GetKey())
		assert.Equal(t, "jonjozwiak/schosterbarak", valueTag.GetValue())
	})

	t.Run("Tag description tests", func(t *testing.T) {
		tag := tags.Tag{}
		defaultDescription := tag.GetDescription()
		cwd, _ := os.Getwd()
		g := TagGroup{}
		g.InitTagGroup(cwd, nil, nil)
		for _, tag := range g.GetTags() {
			assert.NotEqual(t, defaultDescription, tag.GetDescription())
			assert.NotEqual(t, "", tag.GetDescription())
		}
	})

}

func EvaluateTag(t *testing.T, tag tags.ITag, blame gitservice.GitBlame) tags.ITag {
	tag.Init()
	newTag, err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &tags.Tag{}, newTag)

	return newTag
}

func EvaluateTagWithPrefix(t *testing.T, tag tags.ITag, blame gitservice.GitBlame, tagPrefix string) tags.ITag {
	tag.Init()
	tag.SetTagPrefix(tagPrefix)
	newTag, err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
	assert.Equal(t, "", tag.GetValue())
	assert.IsType(t, &tags.Tag{}, newTag)

	return newTag
}
