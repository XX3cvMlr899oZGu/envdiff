package envsummary_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envsummary"
)

func TestCompute_IdenticalMaps(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "1", "B": "2"}
	s := envsummary.Compute(a, b, envsummary.DefaultOptions())

	if s.Total != 2 || s.Equal != 2 || s.Score != 1.0 {
		t.Errorf("expected all equal, got %+v", s)
	}
	if len(s.TopDiffs) != 0 {
		t.Errorf("expected no top diffs, got %v", s.TopDiffs)
	}
}

func TestCompute_MissingKey(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "1"}
	s := envsummary.Compute(a, b, envsummary.DefaultOptions())

	if s.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", s.Missing)
	}
	if s.Score >= 1.0 {
		t.Errorf("score should be < 1.0, got %.2f", s.Score)
	}
}

func TestCompute_ExtraKey(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"A": "1", "B": "2"}
	s := envsummary.Compute(a, b, envsummary.DefaultOptions())

	if s.Extra != 1 {
		t.Errorf("expected 1 extra, got %d", s.Extra)
	}
}

func TestCompute_ChangedKey(t *testing.T) {
	a := map[string]string{"A": "old"}
	b := map[string]string{"A": "new"}
	s := envsummary.Compute(a, b, envsummary.DefaultOptions())

	if s.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", s.Changed)
	}
	if len(s.TopDiffs) == 0 || s.TopDiffs[0] != "A" {
		t.Errorf("expected TopDiffs to contain A, got %v", s.TopDiffs)
	}
}

func TestCompute_EmptyMaps(t *testing.T) {
	s := envsummary.Compute(nil, nil, envsummary.DefaultOptions())
	if s.Score != 1.0 {
		t.Errorf("empty maps should score 1.0, got %.2f", s.Score)
	}
}

func TestCompute_TopNLimit(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5", "F": "6"}
	b := map[string]string{}
	opts := envsummary.Options{TopN: 3}
	s := envsummary.Compute(a, b, opts)

	if len(s.TopDiffs) != 3 {
		t.Errorf("expected 3 top diffs, got %d", len(s.TopDiffs))
	}
}

func TestFormatText_ContainsFields(t *testing.T) {
	a := map[string]string{"X": "1", "Y": "2"}
	b := map[string]string{"X": "changed"}
	s := envsummary.Compute(a, b, envsummary.DefaultOptions())
	out := envsummary.FormatText(s)

	for _, want := range []string{"Total", "Equal", "Missing", "Extra", "Changed", "Score"} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatText missing field %q", want)
		}
	}
}
