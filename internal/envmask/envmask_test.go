package envmask_test

import (
	"testing"

	"github.com/user/envdiff/internal/envmask"
)

func TestApply_MasksSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
	}
	opts := envmask.DefaultOptions()
	out, err := envmask.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
	if out["DB_PASSWORD"] == "supersecret" {
		t.Error("expected DB_PASSWORD to be masked")
	}
	if len(out["DB_PASSWORD"]) == 0 {
		t.Error("expected non-empty masked value")
	}
}

func TestApply_VisibleChars(t *testing.T) {
	env := map[string]string{"API_TOKEN": "abcdef"}
	opts := envmask.DefaultOptions()
	opts.VisibleChars = 3
	out, err := envmask.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !startsWith(out["API_TOKEN"], "abc") {
		t.Errorf("expected prefix 'abc', got %q", out["API_TOKEN"])
	}
}

func TestApply_EmptyValue(t *testing.T) {
	env := map[string]string{"SECRET_KEY": ""}
	opts := envmask.DefaultOptions()
	out, err := envmask.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET_KEY"] == "" {
		t.Error("expected placeholder for empty sensitive value")
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	opts := envmask.DefaultOptions()
	opts.Patterns = []string{`[invalid`}
	_, err := envmask.Apply(env, opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"AUTH_TOKEN": "xyz123"}
	opts := envmask.DefaultOptions()
	opts.Placeholder = "--"
	opts.VisibleChars = 1
	out, err := envmask.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["AUTH_TOKEN"] != "x--" {
		t.Errorf("expected 'x--', got %q", out["AUTH_TOKEN"])
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
