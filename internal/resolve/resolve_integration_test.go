package resolve_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envdiff/internal/parser"
	"github.com/yourusername/envdiff/internal/resolve"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestResolve_FullPipeline_BraceRefs(t *testing.T) {
	path := writeTempEnvFile(t, "BASE_URL=https://example.com\nAPI_URL=${BASE_URL}/api/v1\nHEALTH=${API_URL}/health\n")

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	resolved, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	want := "https://example.com/api/v1/health"
	if resolved["HEALTH"] != want {
		t.Errorf("HEALTH: want %q, got %q", want, resolved["HEALTH"])
	}
}

func TestResolve_FullPipeline_NoRefs(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\nDB_NAME=mydb\n")

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	resolved, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	if len(resolved) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(resolved))
	}
	for k, v := range env {
		if resolved[k] != v {
			t.Errorf("key %s: want %q, got %q", k, v, resolved[k])
		}
	}
}
