package envindex_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envindex"
)

func TestBuild_BasicIndex(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"HOST": "localhost", "PORT": "8080"},
		"dev":  {"HOST": "localhost", "PORT": "9090"},
	}
	entries := envindex.Build(envs, envindex.DefaultOptions())
	if len(entries) == 0 {
		t.Fatal("expected non-empty index")
	}
	found := false
	for _, e := range entries {
		if e.Value == "localhost" {
			found = true
			if len(e.Occurrences) != 2 {
				t.Errorf("expected 2 envs for 'localhost', got %d", len(e.Occurrences))
			}
		}
	}
	if !found {
		t.Error("expected 'localhost' in index")
	}
}

func TestBuild_ExcludeUnique(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "shared", "B": "only-prod"},
		"dev":  {"A": "shared", "C": "only-dev"},
	}
	opts := envindex.DefaultOptions()
	opts.ExcludeUnique = true
	entries := envindex.Build(envs, opts)
	for _, e := range entries {
		if e.Value == "only-prod" || e.Value == "only-dev" {
			t.Errorf("unique value %q should have been excluded", e.Value)
		}
	}
}

func TestBuild_CaseFold(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"MODE": "Production"},
		"dev":  {"MODE": "production"},
	}
	opts := envindex.DefaultOptions()
	opts.CaseFold = true
	entries := envindex.Build(envs, opts)
	for _, e := range entries {
		if e.Value == "production" && len(e.Occurrences) == 2 {
			return
		}
	}
	t.Error("expected case-folded 'production' to appear in both envs")
}

func TestBuild_EmptyEnvs(t *testing.T) {
	entries := envindex.Build(map[string]map[string]string{}, envindex.DefaultOptions())
	if len(entries) != 0 {
		t.Errorf("expected empty index, got %d entries", len(entries))
	}
}

func TestFormatText_ContainsValue(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"HOST": "db.internal"},
	}
	entries := envindex.Build(envs, envindex.DefaultOptions())
	out := envindex.FormatText(entries)
	if !strings.Contains(out, "db.internal") {
		t.Errorf("expected output to contain 'db.internal', got:\n%s", out)
	}
}

func TestFormatText_EmptyEntries(t *testing.T) {
	out := envindex.FormatText(nil)
	if !strings.Contains(out, "empty") {
		t.Errorf("expected '(empty index)' message, got: %s", out)
	}
}
