package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/tagging/tags"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const org = "bridgecrewio"
const repository = "terragoat"
const filePath = "README.md"
const commitHash1 = "47accf06f13b503f3bab06fed7860e72f7523cac"
const commitHash2 = "b2dc884b7439882c4dbe1e660cb1e02a3f84e45d"

func TestTagCreation(t *testing.T) {
	blame := SetupBlame(t)
	t.Run("BcTraceTagCreation", func(t *testing.T) {
		tag := tags.BcTraceTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "BC_TRACE", tag.Key)
		assert.Equal(t, 36, len(tag.Value))
	})

	t.Run("GitOrgTagCreation", func(t *testing.T) {
		tag := tags.GitOrgTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_org", tag.Key)
		assert.Equal(t, org, tag.Value)
	})

	t.Run("GitRepoTagCreation", func(t *testing.T) {
		tag := tags.GitRepoTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_repo", tag.Key)
		assert.Equal(t, repository, tag.Value)
	})

	t.Run("GitFileTagCreation", func(t *testing.T) {
		tag := tags.GitFileTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_file", tag.Key)
		assert.Equal(t, filePath, tag.Value)
	})

	t.Run("GitCommitTagCreation", func(t *testing.T) {
		tag := tags.GitCommitTag{}
		EvaluateTag(t, &tag, blame)
		assert.Equal(t, "git_commit", tag.Key)
		assert.Equal(t, commitHash1, tag.Value)
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
		assert.Equal(t, "schosterbarak/jonjozwiak", tag.Value)
	})

}

func EvaluateTag(t *testing.T, tag tags.ITag, blame git_service.GitBlame) {
	tag.Init()
	err := tag.CalculateValue(&blame)
	if err != nil {
		assert.Fail(t, "Failed to generate BC trace", err)
	}
}

func SetupBlame(t *testing.T) git_service.GitBlame {
	dateStr0 := "2020-03-28T21:42:46.000Z"
	dateStr1 := "2020-03-27T11:56:33.000Z"
	firstCommitDate, err1 := extractDate(dateStr0)
	secondCommitDate, err2 := extractDate(dateStr1)
	if err1 != nil || err2 != nil {
		assert.Fail(t, "Failed to parse static date")
	}
	return git_service.GitBlame{
		GitOrg:        org,
		GitRepository: repository,
		FilePath:      filePath,
		BlamesByLine: map[int]*git.Line{
			0: {
				Author: "schosterbarak@gmail.com",
				Text:   "# Terragoat",
				Date:   firstCommitDate,
				Hash:   plumbing.NewHash(commitHash1),
			},
			1: {
				Author: "jonjozwiak@users.noreply.github.com",
				Text:   "Bridgecrew solution to create vulnerable infrastructure",
				Date:   secondCommitDate,
				Hash:   plumbing.NewHash(commitHash2),
			},
		},
	}
}

func extractDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	parsedDate, err := time.Parse(layout, dateStr)
	return parsedDate, err
}
