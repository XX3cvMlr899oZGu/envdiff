package envset_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/envset"
)

func TestIntersection_CommonKeys(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"B": "2", "C": "9", "D": "4"}
	out := envset.Intersection(a, b)
	if _, ok := out["B"]; !ok {
		t.Error("expected B in intersection")
	}
	if _, ok := out["C"]; !ok {
		t.Error("expected C in intersection")
	}
	if _, ok := out["A"]; ok {
		t.Error("A should not be in intersection")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestIntersection_Empty(t *testing.T) {
	out := envset.Intersection()
	if len(out) != 0 {
		t.Error("expected empty result")
	}
}

func TestUnion_MergesAll(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	c := map[string]string{"A": "overwritten", "C": "3"}
	out := envset.Union(a, b, c)
	if out["A"] != "overwritten" {
		t.Errorf("expected overwritten, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Error("expected B=2")
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}

func TestDifference_RemovesOtherKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	other := map[string]string{"B": "x", "C": "y"}
	out := envset.Difference(base, other)
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
	if out["A"] != "1" {
		t.Error("expected A=1")
	}
}

func TestDifference_NoOthers(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	out := envset.Difference(base)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
