package snapshot_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

// TestSnapshotThenCompare verifies that a saved snapshot can be loaded
// and its env map used directly with the diff.Compare function.
func TestSnapshotThenCompare(t *testing.T) {
	baseEnv := map[string]string{
		"APP_ENV":    "staging",
		"DB_HOST":    "db.staging.internal",
		"SECRET_KEY": "abc123",
	}

	path := tempSnapshotPath(t)
	if err := snapshot.Save(path, ".env.staging", baseEnv); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	currentEnv := map[string]string{
		"APP_ENV":    "production",
		"DB_HOST":    "db.staging.internal",
		"NEW_FEATURE": "enabled",
	}

	results := diff.Compare(snap.Env, currentEnv)

	if len(results) == 0 {
		t.Fatal("expected diff results, got none")
	}

	statuses := make(map[string]string)
	for _, r := range results {
		statuses[r.Key] = r.Status
	}

	if statuses["APP_ENV"] != "mismatch" {
		t.Errorf("expected APP_ENV to be mismatch, got %q", statuses["APP_ENV"])
	}

	if statuses["SECRET_KEY"] != "missing_in_second" {
		t.Errorf("expected SECRET_KEY missing_in_second, got %q", statuses["SECRET_KEY"])
	}

	if statuses["NEW_FEATURE"] != "missing_in_first" {
		t.Errorf("expected NEW_FEATURE missing_in_first, got %q", statuses["NEW_FEATURE"])
	}

	if statuses["DB_HOST"] != "equal" {
		t.Errorf("expected DB_HOST to be equal, got %q", statuses["DB_HOST"])
	}
}
