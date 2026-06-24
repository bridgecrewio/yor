package structure

import (
	"strings"
	"testing"

	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/stretchr/testify/assert"
)

func TestHelmBlock_GetFramework(t *testing.T) {
	b := &HelmBlock{}
	assert.Equal(t, "Helm", b.GetFramework())
}

func TestIsValidK8sLabel(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		valid bool
	}{
		{"yor_trace uuid", "yor_trace", "00000000-0000-0000-0000-000000000001", true},
		{"yor_name simple", "yor_name", "deployment", true},
		{"git_commit hash", "git_commit", "abcdef0123456789", true},
		{"empty value allowed", "yor_name", "", true},
		{"git email invalid char", "git_last_modified_by", "user@example.com", false},
		{"value too long", "git_file", strings.Repeat("a", 64), false},
		{"path with slash in value", "git_file", "src/main/app.yaml", false},
		{"value leading dash", "k", "-bad", false},
		{"prefixed key valid", "app.kubernetes.io/name", "myapp", true},
		{"empty key invalid", "", "x", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsValidK8sLabel(tt.key, tt.value))
		})
	}
}

func TestHelmBlock_GetValidLabels_DropsInvalid(t *testing.T) {
	b := &HelmBlock{
		Block: structure.Block{
			IsTaggable: true,
		},
	}
	b.AddNewTags([]tags.ITag{
		&tags.Tag{Key: "yor_name", Value: "deployment"},
		&tags.Tag{Key: "git_last_modified_by", Value: "user@example.com"}, // invalid (email)
		&tags.Tag{Key: "git_commit", Value: "abcdef0123456789"},
	})
	valid, dropped := b.GetValidLabels()

	validKeys := map[string]bool{}
	for _, v := range valid {
		validKeys[v.GetKey()] = true
	}
	droppedKeys := map[string]bool{}
	for _, d := range dropped {
		droppedKeys[d.GetKey()] = true
	}

	assert.True(t, validKeys["yor_name"])
	assert.True(t, validKeys["git_commit"])
	assert.False(t, validKeys["git_last_modified_by"])
	assert.True(t, droppedKeys["git_last_modified_by"])
}

func TestHelmBlock_UpdateTags_PreservesExistingYorTrace(t *testing.T) {
	existingTrace := "11111111-1111-1111-1111-111111111111"
	b := &HelmBlock{
		Block: structure.Block{
			IsTaggable:  true,
			ExitingTags: []tags.ITag{&tags.Tag{Key: tags.YorTraceTagKey, Value: existingTrace}},
		},
	}
	// Simulate a new run computing a different trace value.
	b.AddNewTags([]tags.ITag{&tags.Tag{Key: tags.YorTraceTagKey, Value: "22222222-2222-2222-2222-222222222222"}})
	b.UpdateTags()

	var traceVal string
	for _, m := range b.MergeTags() {
		if m.GetKey() == tags.YorTraceTagKey {
			traceVal = m.GetValue()
		}
	}
	assert.Equal(t, existingTrace, traceVal, "existing yor_trace must be preserved on re-run (D11)")
}
