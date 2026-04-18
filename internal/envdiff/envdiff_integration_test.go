package envdiff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/envdiff"
)

func TestSummarize_FullPipeline(t *testing.T) {
	a := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
		"DB":   "mydb",
	}
	b := map[string]string{
		"HOST": "remotehost",
		"PORT": "8080",
	}

	results := diff.Compare(a, b)
	stats := envdiff.Summarize(results)

	if stats.Total != 3 {
		t.Errorf("expected 3 total, got %d", stats.Total)
	}
	if stats.Mismatch != 1 {
		t.Errorf("expected 1 mismatch, got %d", stats.Mismatch)
	}
	if stats.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", stats.Missing)
	}
	if stats.Equal != 1 {
		t.Errorf("expected 1 equal, got %d", stats.Equal)
	}
	if !envdiff.HasDifferences(results) {
		t.Error("expected differences to be detected")
	}
}
