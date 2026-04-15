package lint

import (
	"strings"
	"testing"
)

func TestRun_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	issues := Run(env, DefaultOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestRun_EmptyValue(t *testing.T) {
	env := map[string]string{"API_KEY": ""}
	issues := Run(env, DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != Warn {
		t.Errorf("expected WARN, got %s", issues[0].Severity)
	}
	if !strings.Contains(issues[0].Message, "empty") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestRun_KeyWithSpaces(t *testing.T) {
	env := map[string]string{"BAD KEY": "value"}
	issues := Run(env, DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != Error {
		t.Errorf("expected ERROR, got %s", issues[0].Severity)
	}
}

func TestRun_ValueTooLong(t *testing.T) {
	env := map[string]string{"LONG_VAL": strings.Repeat("x", 501)}
	opts := DefaultOptions()
	issues := Run(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "exceeds maximum") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestRun_DisabledChecks(t *testing.T) {
	env := map[string]string{"EMPTY": "", "BAD KEY": "v"}
	opts := Options{WarnEmptyValue: false, WarnKeySpaces: false}
	issues := Run(env, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues when checks disabled, got %d", len(issues))
	}
}

func TestFormatIssues_Empty(t *testing.T) {
	out := FormatIssues(nil)
	if out != "no lint issues found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatIssues_WithIssues(t *testing.T) {
	issues := []Issue{
		{Key: "FOO", Message: "value is empty", Severity: Warn},
	}
	out := FormatIssues(issues)
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected WARN tag in output: %s", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output: %s", out)
	}
}
