package integration

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bridgecrewio/yor/src/common/clioptions"
	"github.com/bridgecrewio/yor/src/common/runner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

// writeChart materializes a minimal Helm chart inside dir and returns the deployment template path.
func writeChart(t *testing.T, dir string) {
	t.Helper()
	mustWrite(t, filepath.Join(dir, "Chart.yaml"), "apiVersion: v2\nname: it-chart\nversion: 0.1.0\n")
	mustWrite(t, filepath.Join(dir, "values.yaml"), "replicaCount: 1\n")
	templates := filepath.Join(dir, "templates")
	failIfErr(t, os.MkdirAll(templates, 0755))
	mustWrite(t, filepath.Join(templates, "_helpers.tpl"), "{{- define \"it.name\" -}}it{{- end -}}\n")
	mustWrite(t, filepath.Join(templates, "NOTES.txt"), "thanks\n")
	mustWrite(t, filepath.Join(templates, "deployment.yaml"),
		"apiVersion: apps/v1\n"+
			"kind: Deployment\n"+
			"metadata:\n"+
			"  name: {{ include \"it.name\" . }}\n"+
			"spec:\n"+
			"  replicas: {{ .Values.replicaCount }}\n"+
			"  selector:\n"+
			"    matchLabels:\n"+
			"      app: it\n"+
			"  template:\n"+
			"    metadata:\n"+
			"      labels:\n"+
			"        app: it\n"+
			"    spec:\n"+
			"      containers:\n"+
			"        - name: app\n"+
			"          image: nginx\n")
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	failIfErr(t, os.WriteFile(path, []byte(content), 0644))
}

func TestHelmIntegration(t *testing.T) {
	t.Run("Helm parser tags workloads with valid labels in a git repo", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "helm-it")
		failIfErr(t, err)
		defer func() { _ = os.RemoveAll(dir) }()

		writeChart(t, dir)

		// Set up a git repo so git_* tags resolve (only valid ones become labels).
		repo, err := git.PlainInit(dir, false)
		failIfErr(t, err)
		worktree, err := repo.Worktree()
		failIfErr(t, err)
		_, err = worktree.Add(".")
		failIfErr(t, err)
		_, err = worktree.Commit("init chart", &git.CommitOptions{
			Author: &object.Signature{Name: "Tester", Email: "tester@example.com", When: time.Now()},
		})
		failIfErr(t, err)

		// Run yor with ONLY the Helm parser (opt-in registration, D10).
		yorRunner := runner.Runner{}
		err = yorRunner.Init(&clioptions.TagOptions{
			Directory: dir,
			TagGroups: getTagGroups(),
			Parsers:   []string{"Helm"},
		})
		failIfErr(t, err)
		_, err = yorRunner.TagDirectory()
		failIfErr(t, err)

		deployment := readFileString(t, path.Join(dir, "templates", "deployment.yaml"))

		// Baseline labels are always written and label-safe (D3.2).
		assert.Contains(t, deployment, "yor_trace:")
		assert.Contains(t, deployment, "yor_name:")
		// Invalid git label values (emails/paths) must NEVER be written as labels (D3.1),
		// regardless of whether git history is rich enough to resolve git_commit.
		assert.NotContains(t, deployment, "git_last_modified_by:")
		assert.NotContains(t, deployment, "git_file:")
		assert.NotContains(t, deployment, "git_last_modified_at:")
		// Template directives are preserved verbatim (D1).
		assert.Contains(t, deployment, `{{ include "it.name" . }}`)

		// Top-level-only label scope (D7): yor tags are written ONLY to the resource's own
		// top-level metadata.labels. One yor_trace occurrence = one (top-level) site.
		assert.Equal(t, 1, strings.Count(deployment, "yor_trace:"), "labels written ONLY to top-level metadata")
		// The pod-template metadata and the immutable selector are NEVER modified (D7).
		templateIdx := strings.Index(deployment, "template:")
		if templateIdx >= 0 {
			nestedBlock := deployment[templateIdx:]
			assert.NotContains(t, nestedBlock, "yor_trace", "spec.template metadata and selector must stay untouched")
			assert.NotContains(t, nestedBlock, "yor_name", "spec.template metadata and selector must stay untouched")
		}

		// Non-resource files must be untouched.
		assert.Equal(t, "thanks\n", readFileString(t, path.Join(dir, "templates", "NOTES.txt")))
		assert.False(t, strings.Contains(readFileString(t, path.Join(dir, "values.yaml")), "yor_trace"))
	})

	t.Run("opt-in respected: without Helm parser the chart is untouched", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "helm-it-optin")
		failIfErr(t, err)
		defer func() { _ = os.RemoveAll(dir) }()

		writeChart(t, dir)

		// git-init so the runner's git tag group does not abort the process when it cannot find a repo.
		repo, err := git.PlainInit(dir, false)
		failIfErr(t, err)
		worktree, err := repo.Worktree()
		failIfErr(t, err)
		_, err = worktree.Add(".")
		failIfErr(t, err)
		_, err = worktree.Commit("init chart", &git.CommitOptions{
			Author: &object.Signature{Name: "Tester", Email: "tester@example.com", When: time.Now()},
		})
		failIfErr(t, err)

		before := readFileString(t, path.Join(dir, "templates", "deployment.yaml"))

		yorRunner := runner.Runner{}
		err = yorRunner.Init(&clioptions.TagOptions{
			Directory: dir,
			TagGroups: getTagGroups(),
			Parsers:   []string{"Terraform", "CloudFormation", "Serverless"},
		})
		failIfErr(t, err)
		_, err = yorRunner.TagDirectory()
		failIfErr(t, err)

		after := readFileString(t, path.Join(dir, "templates", "deployment.yaml"))
		assert.Equal(t, before, after, "Helm chart must be untouched when Helm is not in --parsers (D10)")
	})
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	// #nosec G304 - test file
	b, err := os.ReadFile(path)
	failIfErr(t, err)
	return string(b)
}
