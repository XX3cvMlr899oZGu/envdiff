package envdiff_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/envdiff"
)

var results = []diff.Result{
	{Key: "A", Status: diff.Equal, Value1: "x", Value2: "x"},
	{Key: "B", Status: diff.MissingInSecond, Value1: "y", Value2: ""},
	{Key: "C", Status: diff.MissingInFirst, Value1: "", Value2: "z"},
	{Key: "D", Status: diff.ValueMismatch, Value1: "1", Value2: "2"},
}

func TestSummarize_AllEqual(t *testing.T) {
	only := []diff.Result{{Key: "X", Status: diff.Equal, Value1: "a", Value2: "a"}}
	s := envdiff.Summarize(only)
	if s.Total != 1 || s.Equal != 1 || s.Missing != 0 || s.Extra != 0 || s.Mismatch != 0 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestSummarize_Mixed(t *testing.T) {
	s := envdiff.Summarize(results)
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Equal != 1 {
		t.Errorf("expected Equal=1, got %d", s.Equal)
	}
	if s.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", s.Missing)
	}
	if s.Extra != 1 {
		t.Errorf("expected Extra=1, got %d", s.Extra)
	}
	if s.Mismatch != 1 {
		t.Errorf("expected Mismatch=1, got %d", s.Mismatch)
	}
}

func TestHasDifferences_False(t *testing.T) {
	only := []diff.Result{{Key: "X", Status: diff.Equal}}
	if envdiff.HasDifferences(only) {
		t.Error("expected no differences")
	}
}

func TestHasDifferences_True(t *testing.T) {
	if !envdiff.HasDifferences(results) {
		t.Error("expected differences to be detected")
	}
}

func TestFormatSummary_ContainsFields(t *testing.T) {
	s := envdiff.Summarize(results)
	line := envdiff.FormatSummary(s)
	for _, want := range []string{"total=", "equal=", "missing=", "extra=", "mismatch="} {
		if !contains(line, want) {
			t.Errorf("FormatSummary missing field %q in %q", want, line)
		}
	}
}

func TestKeysByStatus_MissingInSecond(t *testing.T) {
	keys := envdiff.KeysByStatus(results, diff.MissingInSecond)
	if len(keys) != 1 || keys[0] != "B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, sub))
}

func containsRune(s, sub string) bool {
	for i := range s {
		if i+len(sub) <= len(s) && s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
