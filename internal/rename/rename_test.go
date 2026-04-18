package rename

import (
	"strings"
	"testing"
)

func TestApply_NoRenames(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := Apply(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
	if res.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
}

func TestApply_ExplicitMap(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	opts := DefaultOptions()
	opts.Map["OLD_KEY"] = "NEW_KEY"
	res, err := Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY in output")
	}
	if _, ok := res.Env["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if len(res.Changes) != 1 || res.Changes[0].OldKey != "OLD_KEY" {
		t.Errorf("unexpected changes: %+v", res.Changes)
	}
}

func TestApply_PrefixSubstitution(t *testing.T) {
	env := map[string]string{"DEV_HOST": "localhost", "DEV_PORT": "5432", "OTHER": "x"}
	opts := Options{Map: make(map[string]string), OldPrefix: "DEV_", NewPrefix: "PROD_"}
	res, err := Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["PROD_HOST"] != "localhost" {
		t.Errorf("expected PROD_HOST=localhost")
	}
	if res.Env["PROD_PORT"] != "5432" {
		t.Errorf("expected PROD_PORT=5432")
	}
	if res.Env["OTHER"] != "x" {
		t.Errorf("expected OTHER unchanged")
	}
	if len(res.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(res.Changes))
	}
}

func TestApply_ConflictReturnsError(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	opts := DefaultOptions()
	opts.Map["A"] = "B"
	_, err := Apply(env, opts)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

// TestApply_MapPreservesValue ensures that the value associated with a renamed
// key is carried over unchanged to the new key.
func TestApply_MapPreservesValue(t *testing.T) {
	env := map[string]string{"OLD": "secret"}
	opts := DefaultOptions()
	opts.Map["OLD"] = "NEW"
	res, err := Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["NEW"] != "secret" {
		t.Errorf("expected NEW=secret, got %q", res.Env["NEW"])
	}
}

func TestFormatChanges_NoChanges(t *testing.T) {
	out := FormatChanges(nil)
	if !strings.Contains(out, "no keys renamed") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatChanges_WithChanges(t *testing.T) {
	changes := []Change{{OldKey: "FOO", NewKey: "BAR"}}
	out := FormatChanges(changes)
	if !strings.Contains(out, "FOO -> BAR") {
		t.Errorf("expected FOO -> BAR in output, got %q", out)
	}
}
