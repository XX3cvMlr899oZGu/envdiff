package envcheck_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envcheck"
)

func TestCheck_NoViolations(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	rules := []envcheck.Rule{
		{Key: "HOST", Type: envcheck.RuleRequired},
		{Key: "PORT", Type: envcheck.RuleNonEmpty},
	}
	vs := envcheck.Check(env, rules)
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d", len(vs))
	}
}

func TestCheck_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []envcheck.Rule{
		{Key: "HOST", Type: envcheck.RuleRequired},
	}
	vs := envcheck.Check(env, rules)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Key != "HOST" {
		t.Errorf("expected key HOST, got %s", vs[0].Key)
	}
}

func TestCheck_ForbiddenKeyPresent(t *testing.T) {
	env := map[string]string{"DEBUG": "true"}
	rules := []envcheck.Rule{
		{Key: "DEBUG", Type: envcheck.RuleForbidden},
	}
	vs := envcheck.Check(env, rules)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Rule != envcheck.RuleForbidden {
		t.Errorf("expected forbidden rule, got %s", vs[0].Rule)
	}
}

func TestCheck_NonEmptyViolation(t *testing.T) {
	env := map[string]string{"TOKEN": "   "}
	rules := []envcheck.Rule{
		{Key: "TOKEN", Type: envcheck.RuleNonEmpty},
	}
	vs := envcheck.Check(env, rules)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
}

func TestHasViolations_True(t *testing.T) {
	vs := []envcheck.Violation{{Key: "X", Rule: envcheck.RuleRequired, Message: "missing"}}
	if !envcheck.HasViolations(vs) {
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
	if out != "no violations" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatViolations_ContainsRuleType(t *testing.T) {
	vs := []envcheck.Violation{
		{Key: "HOST", Rule: envcheck.RuleRequired, Message: "required key \"HOST\" is missing"},
	}
	out := envcheck.FormatViolations(vs)
	if !strings.Contains(out, "required") {
		t.Errorf("expected output to contain rule type, got: %s", out)
	}
}
