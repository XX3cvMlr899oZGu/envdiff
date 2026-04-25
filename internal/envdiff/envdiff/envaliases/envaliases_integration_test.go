package envaliases_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envaliases"
	"github.com/user/envdiff/internal/loader"
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

func TestEnvAliases_FullPipeline_Rename(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=production\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	opts := envaliases.Options{
		Aliases: map[string]string{
			"DB_HOST": "DATABASE_HOST",
			"DB_PORT": "DATABASE_PORT",
		},
		KeepOriginal: false,
	}

	out, err := envaliases.Apply(env, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if out["DATABASE_PORT"] != "5432" {
		t.Errorf("expected DATABASE_PORT=5432, got %q", out["DATABASE_PORT"])
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", out["APP_ENV"])
	}
}

func TestEnvAliases_FullPipeline_KeepOriginal(t *testing.T) {
	path := writeTempEnvFile(t, "SECRET_KEY=abc123\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	opts := envaliases.Options{
		Aliases:      map[string]string{"SECRET_KEY": "APP_SECRET"},
		KeepOriginal: true,
	}

	out, err := envaliases.Apply(env, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if out["SECRET_KEY"] != "abc123" {
		t.Errorf("expected SECRET_KEY preserved")
	}
	if out["APP_SECRET"] != "abc123" {
		t.Errorf("expected APP_SECRET=abc123")
	}
}
