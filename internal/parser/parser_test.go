package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-test-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseFile_BasicKeyValues(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\nDB_PORT=5432\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := EnvMap{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	for k, v := range expected {
		if env[k] != v {
			t.Errorf("key %q: got %q, want %q", k, env[k], v)
		}
	}
}

func TestParseFile_IgnoresCommentsAndBlankLines(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\n\nKEY=value\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got KEY=%q", env["KEY"])
	}
}

func TestParseFile_StripQuotes(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret"` + "\n" + `TOKEN='abc123'` + "\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret" {
		t.Errorf("expected 'my secret', got %q", env["SECRET"])
	}
	if env["TOKEN"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", env["TOKEN"])
	}
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")

	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
