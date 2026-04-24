package envwatch_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/envdiff/envwatch"
)

// TestEnvWatch_FullPipeline_DiffOnChange verifies that an Event emitted by
// Watch can be fed directly into diff.Compare to produce a meaningful diff.
func TestEnvWatch_FullPipeline_DiffOnChange(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, "HOST=localhost\nPORT=5432\nDEBUG=false\n")

	opts := envwatch.DefaultOptions()
	opts.PollInterval = 40 * time.Millisecond

	done := make(chan struct{})
	defer close(done)

	ch, err := envwatch.Watch(path, opts, done)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	time.Sleep(60 * time.Millisecond)

	newContent := "HOST=localhost\nPORT=5433\nDEBUG=true\nLOG_LEVEL=info\n"
	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var ev envwatch.Event
	select {
	case ev = <-ch:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for watch event")
	}

	results := diff.Compare(ev.OldEnv, ev.NewEnv)

	statuses := make(map[string]string)
	for _, r := range results {
		statuses[r.Key] = r.Status
	}

	if statuses["PORT"] != "mismatch" {
		t.Errorf("PORT status = %q; want %q", statuses["PORT"], "mismatch")
	}
	if statuses["DEBUG"] != "mismatch" {
		t.Errorf("DEBUG status = %q; want %q", statuses["DEBUG"], "mismatch")
	}
	if statuses["LOG_LEVEL"] != "missing_in_first" {
		t.Errorf("LOG_LEVEL status = %q; want %q", statuses["LOG_LEVEL"], "missing_in_first")
	}
	if statuses["HOST"] != "equal" {
		t.Errorf("HOST status = %q; want %q", statuses["HOST"], "equal")
	}
}
