package validate_test

import (
	"testing"

	"github.com/user/envdiff/internal/validate"
)

// TestValidate_FullPipeline simulates loading an env map and running validation
// rules against it, mirroring real CLI usage.
func TestValidate_FullPipeline(t *testing.T) {
	env := map[string]string{
		"APP_ENV":      "production",
		"PORT":         "443",
		"DATABASE_URL": "postgres://prod-host/db",
		"LOG_LEVEL":    "warn",
	}

	rules := []validate.Rule{
		{Key: "APP_ENV", Pattern: `^(development|staging|production)$`, Required: true},
		{Key: "PORT", Pattern: `^\d+$`, Required: true},
		{Key: "DATABASE_URL", Required: true},
		{Key: "SECRET_KEY", Required: true},
		{Key: "LOG_LEVEL", Pattern: `^(debug|info|warn|error)$`},
	}

	violations := validate.ApplyRules(env, rules)

	if len(violations) != 1 {
		t.Fatalf("expected exactly 1 violation, got %d: %+v", len(violations), violations)
	}

	if violations[0].Key != "SECRET_KEY" {
		t.Errorf("expected violation for SECRET_KEY, got %s", violations[0].Key)
	}

	summary := validate.FormatViolations(violations)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestValidate_AllPass(t *testing.T) {
	env := map[string]string{
		"PORT":      "3000",
		"APP_ENV":   "staging",
		"SECRET_KEY": "supersecret",
	}

	rules := []validate.Rule{
		{Key: "PORT", Pattern: `^\d+$`, Required: true},
		{Key: "APP_ENV", Pattern: `^(development|staging|production)$`, Required: true},
		{Key: "SECRET_KEY", Required: true},
	}

	violations := validate.ApplyRules(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got: %+v", violations)
	}
}
