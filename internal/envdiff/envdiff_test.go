package envdiff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/envdiff"
)

func results(statuses ...diff.Status) []diff.Result {
	var rs []diff.Result
	for i, s := range statuses {
		rs = append(rs, diff.Result{Key: fmt.Sprintf("KEY%d", i), Status: s})
	}
	return rs
}

func TestSummarize_AllEqual(t *testing.T) {
	res := []diff.Result{
		{Key: "A", Status: diff.Equal},
		{Key: "B", Status: diff.Equal},
	}
	s := envdiff.Summarize(res)
	if s.Total != 2 || s.Equal != 2 || s.Missing != 0 {
		t.Errorf("unexpected stats: %+v", s)
	}
}

func TestSummarize_Mixed(t *testing.T) {
	res := []diff.Result{
		{Key: "A", Status: diff.Equal},
		{Key: "B", Status: diff.Missing},
		{Key: "C", Status: diff.Mismatch},
		{Key: "D", Status: diff.Extra},
	}
	s := envdiff.Summarize(res)
	if s.Total != 4 || s.Equal != 1 || s.Missing != 1 || s.Mismatch != 1 || s.Extra != 1 {
		t.Errorf("unexpected stats: %+v", s)
	}
}

func TestHasDifferences_False(t *testing.T) {
	res := []diff.Result{{Key: "A", Status: diff.Equal}}
	if envdiff.HasDifferences(res) {
		t.Error("expected no differences")
	}
}

func TestHasDifferences_True(t *testing.T) {
	res := []diff.Result{
		{Key: "A", Status: diff.Equal},
		{Key: "B", Status: diff.Missing},
	}
	if !envdiff.HasDifferences(res) {
		t.Error("expected differences")
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := envdiff.Summarize(nil)
	if s.Total != 0 {
		t.Errorf("expected zero total, got %d", s.Total)
	}
}
