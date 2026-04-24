package envcheck_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/envdiff/envdiff/envcheck"
)

func TestCheck_NoViolations(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	rules := []envcheck.Rule{
		{Key: "HOST", Kind: envcheck.RuleRequired},
		{Key: "PORT", Kind: envcheck.RuleNonEmpty},
	}
	got := envcheck.Check(env, rules)
	if len(got) != 0 {
		t.Errorf("expected no violations, got %d", len(got))
	}
}

func TestCheck_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []envcheck.Rule{
		{Key: "HOST", Kind: envcheck.RuleRequired},
	}
	got := envcheck.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Key != "HOST" {
		t.Errorf("expected violation for HOST, got %s", got[0].Key)
	}
}

func TestCheck_ForbiddenKeyPresent(t *testing.T) {
	env := map[string]string{"DEBUG": "true"}
	rules := []envcheck.Rule{
		{Key: "DEBUG", Kind: envcheck.RuleForbidden},
	}
	got := envcheck.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Kind != envcheck.RuleForbidden {
		t.Errorf("expected forbidden kind, got %s", got[0].Kind)
	}
}

func TestCheck_NonEmptyViolation(t *testing.T) {
	env := map[string]string{"SECRET": "   "}
	rules := []envcheck.Rule{
		{Key: "SECRET", Kind: envcheck.RuleNonEmpty},
	}
	got := envcheck.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
}

func TestHasViolations_True(t *testing.T) {
	v := []envcheck.Violation{{Key: "X", Kind: envcheck.RuleRequired, Message: "missing"}}
	if !envcheck.HasViolations(v) {
		t.Error("expected HasViolations to return true")
	}
}

func TestHasViolations_False(t *testing.T) {
	if envcheck.HasViolations(nil) {
		t.Error("expected HasViolations to return false for nil")
	}
}

func TestFormatViolations_NoViolations(t *testing.T) {
	out := envcheck.FormatViolations(nil)
	if out != "no violations found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatViolations_ContainsKind(t *testing.T) {
	v := []envcheck.Violation{
		{Key: "HOST", Kind: envcheck.RuleRequired, Message: "required key \"HOST\" is missing"},
	}
	out := envcheck.FormatViolations(v)
	if !strings.Contains(out, "[required]") {
		t.Errorf("expected output to contain '[required]', got: %s", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected output to contain 'HOST', got: %s", out)
	}
}
