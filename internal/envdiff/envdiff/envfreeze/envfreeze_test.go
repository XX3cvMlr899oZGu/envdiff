package envfreeze_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envdiff/internal/envdiff/envdiff/envfreeze"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "freeze.json")
}

func TestFreezeAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	path := tempPath(t)

	if err := envfreeze.Freeze(env, path); err != nil {
		t.Fatalf("Freeze: %v", err)
	}

	record, err := envfreeze.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	for k, want := range env {
		if got := record.Env[k]; got != want {
			t.Errorf("key %s: got %q want %q", k, got, want)
		}
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := envfreeze.Load("/nonexistent/freeze.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := envfreeze.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDiff_NoDeviations(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	path := tempPath(t)
	_ = envfreeze.Freeze(env, path)
	record, _ := envfreeze.Load(path)

	devs := envfreeze.Diff(record, env, envfreeze.DefaultOptions())
	if len(devs) != 0 {
		t.Fatalf("expected 0 deviations, got %d", len(devs))
	}
}

func TestDiff_AddedKey(t *testing.T) {
	path := tempPath(t)
	_ = envfreeze.Freeze(map[string]string{"A": "1"}, path)
	record, _ := envfreeze.Load(path)

	devs := envfreeze.Diff(record, map[string]string{"A": "1", "B": "new"}, envfreeze.DefaultOptions())
	if len(devs) != 1 || devs[0].Kind != "added" || devs[0].Key != "B" {
		t.Fatalf("unexpected deviations: %+v", devs)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	path := tempPath(t)
	_ = envfreeze.Freeze(map[string]string{"A": "1", "B": "2"}, path)
	record, _ := envfreeze.Load(path)

	devs := envfreeze.Diff(record, map[string]string{"A": "1"}, envfreeze.DefaultOptions())
	if len(devs) != 1 || devs[0].Kind != "removed" || devs[0].Key != "B" {
		t.Fatalf("unexpected deviations: %+v", devs)
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	path := tempPath(t)
	_ = envfreeze.Freeze(map[string]string{"PORT": "8080"}, path)
	record, _ := envfreeze.Load(path)

	devs := envfreeze.Diff(record, map[string]string{"PORT": "9090"}, envfreeze.DefaultOptions())
	if len(devs) != 1 || devs[0].Kind != "changed" {
		t.Fatalf("unexpected deviations: %+v", devs)
	}
	if devs[0].Frozen != "8080" || devs[0].Current != "9090" {
		t.Errorf("wrong values: %+v", devs[0])
	}
}

func TestDiff_IgnoreKeys(t *testing.T) {
	path := tempPath(t)
	_ = envfreeze.Freeze(map[string]string{"A": "1", "SECRET": "old"}, path)
	record, _ := envfreeze.Load(path)

	opts := envfreeze.DefaultOptions()
	opts.IgnoreKeys = []string{"SECRET"}
	devs := envfreeze.Diff(record, map[string]string{"A": "1", "SECRET": "new"}, opts)
	if len(devs) != 0 {
		t.Fatalf("expected ignored key to be skipped, got %+v", devs)
	}
}

func TestFormatText_NoDeviations(t *testing.T) {
	out := envfreeze.FormatText(nil)
	if out == "" {
		t.Error("expected non-empty output for no deviations")
	}
}

func TestFormatText_ContainsKinds(t *testing.T) {
	devs := []envfreeze.Deviation{
		{Key: "A", Kind: "added", Current: "v"},
		{Key: "B", Kind: "removed", Frozen: "old"},
		{Key: "C", Kind: "changed", Frozen: "x", Current: "y"},
	}
	out := envfreeze.FormatText(devs)
	for _, substr := range []string{"+", "-", "~", "A", "B", "C"} {
		if !containsStr(out, substr) {
			t.Errorf("output missing %q: %s", substr, out)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
