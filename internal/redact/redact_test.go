package redact

import (
	"testing"
)

func TestApply_RedactsSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
		"API_KEY":     "key-abc123",
	}

	result, err := Apply(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "***REDACTED***" {
		t.Errorf("expected API_KEY to be redacted, got %q", result["API_KEY"])
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result["APP_NAME"])
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET_TOKEN": "tok123", "HOST": "localhost"}
	opts := Options{Patterns: DefaultPatterns, Placeholder: "[hidden]"}

	result, err := Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["SECRET_TOKEN"] != "[hidden]" {
		t.Errorf("expected [hidden], got %q", result["SECRET_TOKEN"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST unchanged, got %q", result["HOST"])
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	opts := Options{Patterns: []string{"[invalid"}, Placeholder: "X"}
	_, err := Apply(env, opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	result, err := Apply(map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestIsSensitiveKey(t *testing.T) {
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"AUTH_TOKEN", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"PRIVATE_KEY", true},
	}
	for _, tc := range cases {
		got := IsSensitiveKey(tc.key)
		if got != tc.expected {
			t.Errorf("IsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}
