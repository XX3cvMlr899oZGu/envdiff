package envprune_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/envdiff/envdiff/envprune"
)

func TestApply_NoOptions(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	res, err := envprune.Apply(env, envprune.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals, got %v", res.Removed)
	}
}

func TestApply_RemoveEmpty(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "", "C": "world"}
	opts := envprune.DefaultOptions()
	opts.RemoveEmpty = true
	res, err := envprune.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["B"]; ok {
		t.Error("expected B to be removed")
	}
	if len(res.Removed) != 1 || res.Removed[0] != "B" {
		t.Errorf("expected [B] removed, got %v", res.Removed)
	}
}

func TestApply_RemoveDuplicates(t *testing.T) {
	env := map[string]string{"X": "same", "Y": "same", "Z": "different"}
	opts := envprune.DefaultOptions()
	opts.RemoveDuplicates = true
	res, err := envprune.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys after dedup, got %d", len(res.Env))
	}
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removal, got %v", res.Removed)
	}
}

func TestApply_PatternKeys(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc", "SECRET_TOKEN": "xyz", "HOST": "localhost"}
	opts := envprune.DefaultOptions()
	opts.PatternKeys = "^SECRET_"
	res, err := envprune.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["HOST"]; !ok {
		t.Error("expected HOST to remain")
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %v", res.Removed)
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	env := map[string]string{"A": "1"}
	opts := envprune.DefaultOptions()
	opts.PatternKeys = "[invalid"
	_, err := envprune.Apply(env, opts)
	if err == nil {
		t.Error("expected error for invalid regexp, got nil")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"A": "", "B": "keep"}
	opts := envprune.DefaultOptions()
	opts.RemoveEmpty = true
	envprune.Apply(env, opts) //nolint:errcheck
	if _, ok := env["A"]; !ok {
		t.Error("original map was mutated")
	}
}
