package diff

import (
	"testing"
)

func TestCompare_Nodifferences(t *testing.T) {
	first := map[string]string{"KEY": "value", "FOO": "bar"}
	second := map[string]string{"KEY": "value", "FOO": "bar"}

	result := Compare(first, second)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingInSecond(t *testing.T) {
	first := map[string]string{"KEY": "value", "ONLY_FIRST": "x"}
	second := map[string]string{"KEY": "value"}

	result := Compare(first, second)

	if len(result.MissingInSecond) != 1 || result.MissingInSecond[0] != "ONLY_FIRST" {
		t.Errorf("expected ONLY_FIRST missing in second, got %v", result.MissingInSecond)
	}
	if len(result.MissingInFirst) != 0 {
		t.Errorf("expected no keys missing in first, got %v", result.MissingInFirst)
	}
}

func TestCompare_MissingInFirst(t *testing.T) {
	first := map[string]string{"KEY": "value"}
	second := map[string]string{"KEY": "value", "ONLY_SECOND": "y"}

	result := Compare(first, second)

	if len(result.MissingInFirst) != 1 || result.MissingInFirst[0] != "ONLY_SECOND" {
		t.Errorf("expected ONLY_SECOND missing in first, got %v", result.MissingInFirst)
	}
	if len(result.MissingInSecond) != 0 {
		t.Errorf("expected no keys missing in second, got %v", result.MissingInSecond)
	}
}

func TestCompare_MismatchedValues(t *testing.T) {
	first := map[string]string{"KEY": "old_value"}
	second := map[string]string{"KEY": "new_value"}

	result := Compare(first, second)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatched key, got %d", len(result.Mismatched))
	}
	mm := result.Mismatched[0]
	if mm.Key != "KEY" || mm.First != "old_value" || mm.Second != "new_value" {
		t.Errorf("unexpected mismatch entry: %+v", mm)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := Compare(map[string]string{}, map[string]string{})

	if result.HasDifferences() {
		t.Errorf("expected no differences for two empty maps")
	}
}
