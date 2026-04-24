package envcheck_test

import (
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/envdiff/envdiff/envcheck"
)

func TestCheck_NoViolations(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	rules := []envcheck.Rule{
		{Key: "DB_HOST", Kind: envcheck.Required},
		{Key: "DB_PORT", Kind: envcheck.NonEmpty},
	}
	got := envcheck.Check(env, rules)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %d", len(got))
	}
}

func TestCheck_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{}
	rules := []envcheck.Rule{{Key: "API_KEY", Kind: envcheck.Required}}
	got := envcheck.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Key != "API_KEY" {
		t.Errorf("unexpected key: %s", got[0].Key)
	}
}

func TestCheck_ForbiddenKeyPresent(t *testing.T) {
	env := map[string]string{"DEBUG": "true"}
	rules := []envcheck.Rule{{Key: "DEBUG", Kind: envcheck.Forbidden}}
	got := envcheck.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Kind != envcheck.Forbidden {
		t.Errorf("expected Forbidden kind, got %s", got[0].Kind)
	}
}

func TestCheck_NonEmptyViolation(t *testing.T) {
	env := map[string]string{"SECRET": "   "}
	rules := []envcheck.Rule{{Key: "SECRET", Kind: envcheck.NonEmpty}}
	got := envcheck.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
}

func TestHasViolations_True(t *testing.T) {
	v := []envcheck.Violation{{Key: "X", Kind: envcheck.Required, Message: "missing"}}
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
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatViolations_ContainsMessages(t *testing.T) {
	v := []envcheck.Violation{
		{Key: "DB_HOST", Kind: envcheck.Required, Message: `required key "DB_HOST" is missing`},
		{Key: "DEBUG", Kind: envcheck.Forbidden, Message: `forbidden key "DEBUG" is present`},
	}
	out := envcheck.FormatViolations(v)
	if !strings.Contains(out, "required") {
		t.Error("expected output to contain 'required'")
	}
	if !strings.Contains(out, "forbidden") {
		t.Error("expected output to contain 'forbidden'")
	}
}
