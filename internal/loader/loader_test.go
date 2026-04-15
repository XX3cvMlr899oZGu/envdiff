package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" || env["BAZ"] != "qux" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := loader.LoadFile("/nonexistent/path/.env", loader.Options{})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFile_WithPrefixFilter(t *testing.T) {
	path := writeTempEnv(t, "APP_HOST=localhost\nDB_HOST=db\nAPP_PORT=8080\n")
	env, err := loader.LoadFile(path, loader.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 2 {
		t.Errorf("expected 2 keys after prefix filter, got %d", len(env))
	}
	if _, ok := env["DB_HOST"]; ok {
		t.Error("DB_HOST should have been filtered out")
	}
}

func TestLoadFile_WithExclude(t *testing.T) {
	path := writeTempEnv(t, "SECRET=abc\nPUBLIC=xyz\n")
	env, err := loader.LoadFile(path, loader.Options{Exclude: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["SECRET"]; ok {
		t.Error("SECRET should have been excluded")
	}
	if env["PUBLIC"] != "xyz" {
		t.Errorf("expected PUBLIC=xyz, got %s", env["PUBLIC"])
	}
}

func TestLoadFiles_MultipleFiles(t *testing.T) {
	path1 := writeTempEnv(t, "A=1\n")
	path2 := writeTempEnv(t, "B=2\n")
	results, err := loader.LoadFiles([]string{path1, path2}, loader.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Env["A"] != "1" {
		t.Errorf("expected A=1 in first file")
	}
	if results[1].Env["B"] != "2" {
		t.Errorf("expected B=2 in second file")
	}
}
