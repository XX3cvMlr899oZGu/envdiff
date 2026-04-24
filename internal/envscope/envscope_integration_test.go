package envscope_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/envscope"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/filter"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestEnvScope_FullPipeline_StripPrefix(t *testing.T) {
	path := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=postgres://localhost/db\n")

	env, err := loader.LoadFile(path, filter.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	scope := envscope.Scope{Name: "app", Prefix: "APP_"}
	opts := envscope.Options{StripPrefix: true}
	result, err := envscope.Apply(env, scope, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", result["PORT"])
	}
	if _, ok := result["DB_URL"]; ok {
		t.Error("DB_URL should not be in app scope")
	}
}

func TestEnvScope_FullPipeline_MultiScope(t *testing.T) {
	path := writeTempEnvFile(t, "APP_HOST=localhost\nDB_HOST=db.local\nLOG_LEVEL=debug\n")

	env, err := loader.LoadFile(path, filter.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	scopes := []envscope.Scope{
		{Name: "app", Prefix: "APP_"},
		{Name: "db", Prefix: "DB_"},
	}
	out, err := envscope.ApplyAll(env, scopes, envscope.DefaultOptions())
	if err != nil {
		t.Fatalf("ApplyAll: %v", err)
	}
	if len(out["app"]) != 1 {
		t.Errorf("expected 1 app key, got %d", len(out["app"]))
	}
	if len(out["db"]) != 1 {
		t.Errorf("expected 1 db key, got %d", len(out["db"]))
	}
}
