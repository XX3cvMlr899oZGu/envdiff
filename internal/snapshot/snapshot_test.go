package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/snapshot"
)

func tempSnapshotPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "snap.json")
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
	}
	path := tempSnapshotPath(t)

	if err := snapshot.Save(path, ".env.production", env); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if snap.Source != ".env.production" {
		t.Errorf("expected source %q, got %q", ".env.production", snap.Source)
	}

	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}

	if snap.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}

	for k, v := range env {
		if got := snap.Env[k]; got != v {
			t.Errorf("key %q: expected %q, got %q", k, v, got)
		}
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not valid json{"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSave_EmptyEnv(t *testing.T) {
	path := tempSnapshotPath(t)
	if err := snapshot.Save(path, ".env", map[string]string{}); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(snap.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(snap.Env))
	}
}
