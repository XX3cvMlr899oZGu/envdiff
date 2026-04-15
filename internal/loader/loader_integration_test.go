package loader_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/loader"
)

// TestLoaderThenCompare verifies that loading two env files and comparing
// them produces the expected diff results end-to-end.
func TestLoaderThenCompare(t *testing.T) {
	path1 := writeTempEnv(t, "HOST=localhost\nPORT=8080\nSECRET=abc\n")
	path2 := writeTempEnv(t, "HOST=localhost\nPORT=9090\nNEW_KEY=value\n")

	opts := loader.Options{Exclude: []string{"SECRET"}}

	envs, err := loader.LoadFiles([]string{path1, path2}, opts)
	if err != nil {
		t.Fatalf("LoadFiles failed: %v", err)
	}

	if len(envs) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(envs))
	}

	// Verify SECRET was excluded from both
	for _, ne := range envs {
		if _, ok := ne.Env["SECRET"]; ok {
			t.Errorf("%s: SECRET should have been excluded", filepath.Base(ne.Name))
		}
	}

	results := diff.Compare(envs[0].Env, envs[1].Env)

	statuses := make(map[string]string)
	for _, r := range results {
		statuses[r.Key] = r.Status
	}

	if statuses["HOST"] != "equal" {
		t.Errorf("expected HOST to be equal, got %s", statuses["HOST"])
	}
	if statuses["PORT"] != "mismatch" {
		t.Errorf("expected PORT to be mismatch, got %s", statuses["PORT"])
	}
	if statuses["NEW_KEY"] != "missing_in_first" {
		t.Errorf("expected NEW_KEY missing_in_first, got %s", statuses["NEW_KEY"])
	}
}
