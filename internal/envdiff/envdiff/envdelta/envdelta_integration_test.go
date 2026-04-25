package envdelta_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envdelta"
	"github.com/user/envdiff/internal/loader"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestEnvDelta_FullPipeline_DetectsChanges(t *testing.T) {
	baseFile := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\nAPP_DEBUG=true\n")
	nextFile := writeTempEnvFile(t, "APP_HOST=prod.example.com\nAPP_PORT=8080\nAPP_NEWKEY=hello\n")

	baseEnv, err := loader.LoadFile(baseFile, loader.Options{})
	if err != nil {
		t.Fatalf("load base: %v", err)
	}
	nextEnv, err := loader.LoadFile(nextFile, loader.Options{})
	if err != nil {
		t.Fatalf("load next: %v", err)
	}

	delta := envdelta.Compute(baseEnv, nextEnv, envdelta.DefaultOptions())

	if !delta.HasChanges() {
		t.Fatal("expected changes")
	}

	changed := delta.ByStatus(envdelta.StatusChanged)
	if len(changed) != 1 || changed[0].Key != "APP_HOST" {
		t.Fatalf("expected APP_HOST changed, got %v", changed)
	}

	removed := delta.ByStatus(envdelta.StatusRemoved)
	if len(removed) != 1 || removed[0].Key != "APP_DEBUG" {
		t.Fatalf("expected APP_DEBUG removed, got %v", removed)
	}

	added := delta.ByStatus(envdelta.StatusAdded)
	if len(added) != 1 || added[0].Key != "APP_NEWKEY" {
		t.Fatalf("expected APP_NEWKEY added, got %v", added)
	}
}

func TestEnvDelta_FullPipeline_NoChanges(t *testing.T) {
	content := "DB_HOST=localhost\nDB_PORT=5432\n"
	baseFile := writeTempEnvFile(t, content)
	nextFile := writeTempEnvFile(t, content)

	baseEnv, _ := loader.LoadFile(baseFile, loader.Options{})
	nextEnv, _ := loader.LoadFile(nextFile, loader.Options{})

	delta := envdelta.Compute(baseEnv, nextEnv, envdelta.DefaultOptions())
	if delta.HasChanges() {
		t.Fatal("expected no changes for identical files")
	}
}
