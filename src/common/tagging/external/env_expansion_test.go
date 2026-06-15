package external

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEvaluateTemplateVariable_F001 covers the F-001 hardening of the
// ${env:NAME} expansion performed against values loaded from --config-file.
//
// Background: an attacker who can edit the external tag-group config file
// (commonly committed alongside the IaC it tags) could previously reference
// arbitrary environment variables — including CI secrets — and have Yor write
// the expanded value into the IaC files (which are then committed and pushed
// by entrypoint.sh in the GitHub Action default flow).
//
// The hardening introduces a strict allowlist (built-in + opt-in via
// YOR_ENV_ALLOWLIST), a hard denylist for well-known secret-bearing names,
// and a global kill switch (YOR_DISABLE_ENV_EXPANSION).
func TestEvaluateTemplateVariable_F001(t *testing.T) {
	t.Run("backward-compatible: GIT_BRANCH is allowed by default", func(t *testing.T) {
		t.Setenv("GIT_BRANCH", "main")
		assert.Equal(t, "main", evaluateTemplateVariable("${env:GIT_BRANCH}"))
	})

	t.Run("backward-compatible: YOR_-prefixed names are allowed by default", func(t *testing.T) {
		t.Setenv("YOR_OWNER", "platform-team")
		assert.Equal(t, "platform-team", evaluateTemplateVariable("${env:YOR_OWNER}"))
	})

	t.Run("non-allowlisted name is NOT expanded (literal returned)", func(t *testing.T) {
		// This is the exact F-001 attack pattern: an arbitrary CI secret
		// referenced from an attacker-controlled config file. Yor must leave
		// the placeholder verbatim rather than substitute the secret.
		t.Setenv("CI_RANDOM_CONFIG", "some-value")
		got := evaluateTemplateVariable("${env:CI_RANDOM_CONFIG}")
		assert.Equal(t, "${env:CI_RANDOM_CONFIG}", got,
			"non-allowlisted env vars must not be substituted into IaC tags")
	})

	t.Run("GITHUB_TOKEN is denied even if explicitly added to allowlist", func(t *testing.T) {
		// Demonstrates that the denylist wins over the allowlist — operator
		// misconfiguration cannot re-open F-001 for well-known secrets.
		t.Setenv("GITHUB_TOKEN", "ghp_supersecret_token_value")
		t.Setenv("YOR_ENV_ALLOWLIST", "GITHUB_TOKEN")
		got := evaluateTemplateVariable("${env:GITHUB_TOKEN}")
		assert.Equal(t, "${env:GITHUB_TOKEN}", got)
		assert.NotContains(t, got, "ghp_supersecret_token_value")
	})

	t.Run("AWS credentials are denied", func(t *testing.T) {
		for _, name := range []string{
			"AWS_SECRET_ACCESS_KEY",
			"AWS_SESSION_TOKEN",
			"AWS_ACCESS_KEY_ID",
		} {
			t.Setenv(name, "AKIAIOSFODNN7EXAMPLE")
			t.Setenv("YOR_ENV_ALLOWLIST", name)
			got := evaluateTemplateVariable("${env:" + name + "}")
			assert.Equal(t, "${env:"+name+"}", got, "must not expand %s", name)
		}
	})

	t.Run("substring denylist catches credential-bearing custom names", func(t *testing.T) {
		// Operator tries to allowlist a custom var whose name looks credential-
		// bearing — denylist substrings (SECRET, PASSWORD, ...) reject it.
		for _, name := range []string{
			"MY_API_SECRET",
			"DEPLOY_PASSWORD",
			"DB_PASSWD",
			"SERVICE_PRIVATE_KEY",
			"VAULT_APIKEY",
			"OAUTH_API_KEY",
			"VENDOR_CREDENTIAL",
		} {
			t.Setenv(name, "leaked-value-"+name)
			t.Setenv("YOR_ENV_ALLOWLIST", name)
			got := evaluateTemplateVariable("${env:" + name + "}")
			assert.Equal(t, "${env:"+name+"}", got, "denylist substring should reject %s", name)
		}
	})

	t.Run("YOR_ENV_ALLOWLIST extends the allowlist for benign names", func(t *testing.T) {
		t.Setenv("DEPLOY_REGION", "us-east-1")
		t.Setenv("YOR_ENV_ALLOWLIST", "DEPLOY_REGION,SOME_OTHER")
		assert.Equal(t, "us-east-1", evaluateTemplateVariable("${env:DEPLOY_REGION}"))
	})

	t.Run("YOR_ENV_ALLOWLIST supports prefix glob form", func(t *testing.T) {
		t.Setenv("MYORG_ENV", "prod")
		t.Setenv("YOR_ENV_ALLOWLIST", "MYORG_*")
		assert.Equal(t, "prod", evaluateTemplateVariable("${env:MYORG_ENV}"))
	})

	t.Run("kill switch disables all expansion, including allowlisted names", func(t *testing.T) {
		t.Setenv("GIT_BRANCH", "main")
		t.Setenv("YOR_DISABLE_ENV_EXPANSION", "1")
		assert.Equal(t, "${env:GIT_BRANCH}", evaluateTemplateVariable("${env:GIT_BRANCH}"))
	})

	t.Run("non-template input is returned unchanged", func(t *testing.T) {
		assert.Equal(t, "literal-value", evaluateTemplateVariable("literal-value"))
	})

	t.Run("missing allowlisted var returns the literal placeholder (no panic)", func(t *testing.T) {
		// Ensure we don't substitute an empty string when the allowlisted var
		// happens to be unset — leaving the placeholder is the safer behavior
		// and matches the pre-existing semantics for missing vars.
		_ = "GIT_BRANCH" // documented allowlisted name
		// Deliberately do NOT set GIT_BRANCH for this subtest scope — but
		// other subtests run in the same process; use a unique allowlisted
		// name we control instead.
		t.Setenv("YOR_ENV_ALLOWLIST", "YOR_F001_UNSET_PROBE")
		got := evaluateTemplateVariable("${env:YOR_F001_UNSET_PROBE}")
		assert.Equal(t, "${env:YOR_F001_UNSET_PROBE}", got)
	})
}

func TestIsEnvVarExpansionAllowed(t *testing.T) {
	cases := []struct {
		name    string
		envName string
		setup   func(t *testing.T)
		allowed bool
	}{
		{"GIT_BRANCH allowed", "GIT_BRANCH", nil, true},
		{"YOR_ prefix allowed", "YOR_FOO", nil, true},
		{"unrelated name denied", "PATH", nil, false},
		{"GITHUB_TOKEN always denied", "GITHUB_TOKEN", nil, false},
		{"denylist substring SECRET", "MY_SECRET_X", nil, false},
		{"denylist substring PASSWORD", "ADMIN_PASSWORD", nil, false},
		{
			"allowlist extension via env",
			"DEPLOY_REGION",
			func(t *testing.T) { t.Setenv("YOR_ENV_ALLOWLIST", "DEPLOY_REGION") },
			true,
		},
		{
			"allowlist extension cannot bypass denylist",
			"GITHUB_TOKEN",
			func(t *testing.T) { t.Setenv("YOR_ENV_ALLOWLIST", "GITHUB_TOKEN") },
			false,
		},
		{
			"allowlist glob matches",
			"ACME_FOO",
			func(t *testing.T) { t.Setenv("YOR_ENV_ALLOWLIST", "ACME_*") },
			true,
		},
		{
			"allowlist glob still subject to denylist",
			"ACME_SECRET",
			func(t *testing.T) { t.Setenv("YOR_ENV_ALLOWLIST", "ACME_*") },
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t)
			}
			assert.Equal(t, tc.allowed, isEnvVarExpansionAllowed(tc.envName))
		})
	}
}

func TestIsTruthy(t *testing.T) {
	for _, v := range []string{"1", "true", "TRUE", "Yes", "on", " on "} {
		assert.True(t, isTruthy(v), "expected truthy: %q", v)
	}
	for _, v := range []string{"", "0", "false", "no", "off", "maybe"} {
		assert.False(t, isTruthy(v), "expected falsy: %q", v)
	}
}
