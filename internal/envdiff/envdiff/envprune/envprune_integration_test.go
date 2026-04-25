package envprune_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envdiff/internal/envdiff/envdiff/envprune"
	"github.com/yourusername/envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestEnvPrune_FullPipeline_RemoveEmpty(t *testing.T) {
	path := writeTempEnvFile(t, "HOST=localhost\nPORT=\nDEBUG=true\n")

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	opts := envprune.DefaultOptions()
	opts.RemoveEmpty = true
	res, err := envprune.Apply(env, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if _, ok := res.Env["PORT"]; ok {
		t.Error("expected PORT to be pruned")
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", res.Env["HOST"])
	}
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removed key, got %d", len(res.Removed))
	}
}

func TestEnvPrune_FullPipeline_PatternKeys(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=db\nDB_PASS=secret\nAPP_NAME=myapp\n")

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	opts := envprune.DefaultOptions()
	opts.PatternKeys = "^DB_"
	res, err := envprune.Apply(env, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if len(res.Env) != 1 {
		t.Errorf("expected 1 key remaining, got %d", len(res.Env))
	}
	if res.Env["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", res.Env["APP_NAME"])
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed keys, got %d", len(res.Removed))
	}
}
