package validate_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/validate"
)

var baseEnv = map[string]string{
	"DATABASE_URL": "postgres://localhost:5432/db",
	"PORT":         "8080",
	"DEBUG":        "true",
}

func TestApplyRules_NoViolations(t *testing.T) {
	rules := []validate.Rule{
		{Key: "PORT", Pattern: `^\d+$`, Required: true},
		{Key: "DEBUG", Pattern: `^(true|false)$`},
	}
	violations := validate.ApplyRules(baseEnv, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestApplyRules_RequiredKeyMissing(t *testing.T) {
	rules := []validate.Rule{
		{Key: "SECRET_KEY", Required: true},
	}
	violations := validate.ApplyRules(baseEnv, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "SECRET_KEY" {
		t.Errorf("expected violation for SECRET_KEY, got %s", violations[0].Key)
	}
}

func TestApplyRules_PatternMismatch(t *testing.T) {
	rules := []validate.Rule{
		{Key: "PORT", Pattern: `^[a-z]+$`},
	}
	violations := validate.ApplyRules(baseEnv, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "does not match pattern") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestApplyRules_OptionalMissingKeySkipped(t *testing.T) {
	rules := []validate.Rule{
		{Key: "OPTIONAL_KEY", Pattern: `^\d+$`, Required: false},
	}
	violations := validate.ApplyRules(baseEnv, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations for optional missing key, got %d", len(violations))
	}
}

func TestFormatViolations_NoViolations(t *testing.T) {
	out := validate.FormatViolations(nil)
	if !strings.Contains(out, "no validation violations") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatViolations_WithViolations(t *testing.T) {
	violations := []validate.Violation{
		{Key: "PORT", Message: "value \"abc\" does not match pattern"},
	}
	out := validate.FormatViolations(violations)
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
	if !strings.Contains(out, "[VIOLATION]") {
		t.Errorf("expected [VIOLATION] tag in output, got: %s", out)
	}
}
