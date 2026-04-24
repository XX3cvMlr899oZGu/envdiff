package envdiff_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/envdiff"
)

func TestSummarize_FullPipeline(t *testing.T) {
	a := map[string]string{
		"HOST":  "localhost",
		"PORT":  "5432",
		"DEBUG": "true",
	}
	b := map[string]string{
		"HOST":    "prod.example.com",
		"PORT":    "5432",
		"TIMEOUT": "30s",
	}

	results := diff.Compare(a, b)

	s := envdiff.Summarize(results)
	if s.Total == 0 {
		t.Fatal("expected non-zero total")
	}
	if !envdiff.HasDifferences(results) {
		t.Error("expected differences between maps")
	}

	// HOST differs in value
	mismatched := envdiff.KeysByStatus(results, diff.ValueMismatch)
	if len(mismatched) != 1 || mismatched[0] != "HOST" {
		t.Errorf("expected HOST as mismatch, got %v", mismatched)
	}

	// DEBUG missing in second
	missing := envdiff.KeysByStatus(results, diff.MissingInSecond)
	if len(missing) != 1 || missing[0] != "DEBUG" {
		t.Errorf("expected DEBUG as missing in second, got %v", missing)
	}

	// TIMEOUT extra (missing in first)
	extra := envdiff.KeysByStatus(results, diff.MissingInFirst)
	if len(extra) != 1 || extra[0] != "TIMEOUT" {
		t.Errorf("expected TIMEOUT as extra, got %v", extra)
	}

	line := envdiff.FormatSummary(s)
	if line == "" {
		t.Error("FormatSummary returned empty string")
	}
}
