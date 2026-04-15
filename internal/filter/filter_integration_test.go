package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/filter"
)

// TestFilterThenCompare verifies that filtering is applied before comparison,
// so excluded / prefix-restricted keys do not appear in diff results.
func TestFilterThenCompare(t *testing.T) {
	env1 := map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_PASSWORD": "secret1",
	}
	env2 := map[string]string{
		"APP_HOST":    "prod.example.com",
		"APP_PORT":    "8080",
		"DB_PASSWORD": "secret2",
	}

	opts := filter.Options{
		Prefix:  "APP_",
	}

	f1, err := filter.ApplyToMap(env1, opts)
	if err != nil {
		t.Fatalf("filter env1: %v", err)
	}
	f2, err := filter.ApplyToMap(env2, opts)
	if err != nil {
		t.Fatalf("filter env2: %v", err)
	}

	results := diff.Compare(f1, f2)

	// DB_PASSWORD must not appear in any diff result
	for _, r := range results {
		if r.Key == "DB_PASSWORD" {
			t.Error("DB_PASSWORD should have been filtered out before comparison")
		}
	}

	// APP_HOST should be a mismatch
	found := false
	for _, r := range results {
		if r.Key == "APP_HOST" {
			found = true
			if r.Status != "mismatch" {
				t.Errorf("expected mismatch for APP_HOST, got %s", r.Status)
			}
		}
	}
	if !found {
		t.Error("expected APP_HOST in diff results")
	}
}
