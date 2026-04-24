package envcheck_test

import (
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/envdiff/envdiff/envcheck"
)

var sampleEnv = map[string]string{
	"APP_NAME":  "myapp",
	"APP_PORT":  "8080",
	"DEBUG":     "",
	"SECRET_KEY": "abc123",
}

func TestCheck_NoViolations(t *testing.T) {
	rules := []envcheck.Rule{
		{Key: "APP_NAME", Required: true, NonEmpty: true},
		{Key: "APP_PORT", Required: true},
	}
	violations := envcheck.Check(sampleEnv, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestCheck_RequiredKeyMissing(t *testing.T) {
	rules := []envcheck.Rule{
		{Key: "DATABASE_URL", Required: true},
	}
	violations := envcheck.Check(sampleEnv, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected key: %s", violations[0].Key)
	}
}

func TestCheck_ForbiddenKeyPresent(t *testing.T) {
	rules := []envcheck.Rule{
		{Key: "SECRET_KEY", Forbidden: true},
	}
	violations := envcheck.Check(sampleEnv, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "forbidden") {
		t.Errorf("expected 'forbidden' in message, got: %s", violations[0].Message)
	}
}

func TestCheck_NonEmptyViolation(t *testing.T) {
	rules := []envcheck.Rule{
		{Key: "DEBUG", NonEmpty: true},
	}
	violations := envcheck.Check(sampleEnv, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "empty") {
		t.Errorf("expected 'empty' in message, got: %s", violations[0].Message)
	}
}

func TestHasViolations_True(t *testing.T) {
	rules := []envcheck.Rule{{Key: "MISSING", Required: true}}
	if !envcheck.HasViolations(sampleEnv, rules) {
		t.Error("expected HasViolations to return true")
	}
}

func TestHasViolations_False(t *testing.T) {
	rules := []envcheck.Rule{{Key: "APP_NAME", Required: true}}
	if envcheck.HasViolations(sampleEnv, rules) {
		t.Error("expected HasViolations to return false")
	}
}

func TestFormatViolations_NoViolations(t *testing.T) {
	out := envcheck.FormatViolations(nil)
	if out != "no violations found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatViolations_WithViolations(t *testing.T) {
	v := []envcheck.Violation{
		{Key: "FOO", Message: "required key is missing"},
	}
	out := envcheck.FormatViolations(v)
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output, got: %s", out)
	}
	if !strings.Contains(out, "1 violation") {
		t.Errorf("expected violation count in output, got: %s", out)
	}
}
