package envfreeze_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envdiff/internal/envdiff/envdiff/envfreeze"
	"github.com/yourusername/envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return path
}

func TestEnvFreeze_FullPipeline_NoDrift(t *testing.T) {
	envPath := writeTempEnvFile(t, "APP_ENV=staging\nPORT=3000\n")
	freezePath := filepath.Join(t.TempDir(), "freeze.json")

	env, err := parser.ParseFile(envPath)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	if err := envfreeze.Freeze(env, freezePath); err != nil {
		t.Fatalf("Freeze: %v", err)
	}

	record, err := envfreeze.Load(freezePath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	devs := envfreeze.Diff(record, env, envfreeze.DefaultOptions())
	if envfreeze.HasDeviations(devs) {
		t.Fatalf("expected no deviations, got: %s", envfreeze.FormatText(devs))
	}
}

func TestEnvFreeze_FullPipeline_DetectsDeviation(t *testing.T) {
	origPath := writeTempEnvFile(t, "APP_ENV=staging\nPORT=3000\nDEBUG=false\n")
	freezePath := filepath.Join(t.TempDir(), "freeze.json")

	orig, err := parser.ParseFile(origPath)
	if err != nil {
		t.Fatalf("ParseFile original: %v", err)
	}
	if err := envfreeze.Freeze(orig, freezePath); err != nil {
		t.Fatalf("Freeze: %v", err)
	}

	// Simulate a drifted env: PORT changed, DEBUG removed, NEW_KEY added
	driftedPath := writeTempEnvFile(t, "APP_ENV=staging\nPORT=9000\nNEW_KEY=surprise\n")
	drifted, err := parser.ParseFile(driftedPath)
	if err != nil {
		t.Fatalf("ParseFile drifted: %v", err)
	}

	record, err := envfreeze.Load(freezePath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	devs := envfreeze.Diff(record, drifted, envfreeze.DefaultOptions())
	if !envfreeze.HasDeviations(devs) {
		t.Fatal("expected deviations but got none")
	}

	kinds := map[string]bool{}
	for _, d := range devs {
		kinds[d.Kind] = true
	}
	for _, expected := range []string{"added", "removed", "changed"} {
		if !kinds[expected] {
			t.Errorf("expected deviation kind %q not found in: %+v", expected, devs)
		}
	}

	text := envfreeze.FormatText(devs)
	if text == "" {
		t.Error("FormatText returned empty string")
	}
}
