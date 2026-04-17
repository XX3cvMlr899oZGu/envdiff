package envsort

import (
	"testing"
)

func TestApply_AscendingOrder(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	entries := Apply(env, DefaultOptions())
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("pos %d: got %s, want %s", i, e.Key, expected[i])
		}
	}
}

func TestApply_DescendingOrder(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	opts := Options{Order: Descending}
	entries := Apply(env, opts)
	expected := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("pos %d: got %s, want %s", i, e.Key, expected[i])
		}
	}
}

func TestApply_PrefixFirst(t *testing.T) {
	env := map[string]string{"DB_HOST": "a", "APP_NAME": "b", "DB_PORT": "c", "ZEBRA": "d"}
	opts := Options{Order: Ascending, Prefix: "DB_"}
	entries := Apply(env, opts)
	if entries[0].Key != "DB_HOST" && entries[0].Key != "DB_PORT" {
		t.Errorf("expected DB_ prefix key first, got %s", entries[0].Key)
	}
	if entries[1].Key != "DB_HOST" && entries[1].Key != "DB_PORT" {
		t.Errorf("expected DB_ prefix key second, got %s", entries[1].Key)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	entries := Apply(map[string]string{}, DefaultOptions())
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(entries))
	}
}

func TestApply_ValuesPreserved(t *testing.T) {
	env := map[string]string{"KEY": "value123"}
	entries := Apply(env, DefaultOptions())
	if entries[0].Value != "value123" {
		t.Errorf("value mismatch: got %s", entries[0].Value)
	}
}
