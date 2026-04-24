package envwatch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/envdiff/envwatch"
)

func writeTempEnv(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatch_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, "FOO=bar\nBAZ=qux\n")

	opts := envwatch.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	done := make(chan struct{})
	defer close(done)

	ch, err := envwatch.Watch(path, opts, done)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	// Give the watcher time to settle, then mutate the file.
	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(path, []byte("FOO=changed\nBAZ=qux\nNEW=1\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.OldEnv["FOO"] != "bar" {
			t.Errorf("OldEnv[FOO] = %q; want %q", ev.OldEnv["FOO"], "bar")
		}
		if ev.NewEnv["FOO"] != "changed" {
			t.Errorf("NewEnv[FOO] = %q; want %q", ev.NewEnv["FOO"], "changed")
		}
		if _, ok := ev.NewEnv["NEW"]; !ok {
			t.Error("expected NEW key in NewEnv")
		}
		if ev.Path != path {
			t.Errorf("Path = %q; want %q", ev.Path, path)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, "KEY=value\n")

	opts := envwatch.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	done := make(chan struct{})
	defer close(done)

	ch, err := envwatch.Watch(path, opts, done)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev := <-ch:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(300 * time.Millisecond):
		// expected: no event
	}
}

func TestWatch_FileNotFound(t *testing.T) {
	opts := envwatch.DefaultOptions()
	done := make(chan struct{})
	defer close(done)

	_, err := envwatch.Watch("/nonexistent/.env", opts, done)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
