package envhash_test

import (
	"testing"

	"github.com/user/envdiff/internal/envhash"
)

func TestCompute_DeterministicOutput(t *testing.T) {
	env := map[string]string{"B": "2", "A": "1", "C": "3"}
	opts := envhash.DefaultOptions()

	h1, err := envhash.Compute(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, err := envhash.Compute(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h1 != h2 {
		t.Errorf("expected identical hashes, got %q and %q", h1, h2)
	}
}

func TestCompute_DifferentEnvsProduceDifferentHashes(t *testing.T) {
	a := map[string]string{"KEY": "val1"}
	b := map[string]string{"KEY": "val2"}
	opts := envhash.DefaultOptions()

	ha, _ := envhash.Compute(a, opts)
	hb, _ := envhash.Compute(b, opts)
	if ha == hb {
		t.Error("expected different hashes for different envs")
	}
}

func TestCompute_ExcludeKeys(t *testing.T) {
	env := map[string]string{"A": "1", "SECRET": "topsecret"}
	opts := envhash.Options{ExcludeKeys: []string{"SECRET"}}

	h, err := envhash.Compute(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	envNoSecret := map[string]string{"A": "1"}
	hNoSecret, _ := envhash.Compute(envNoSecret, envhash.DefaultOptions())

	if h != hNoSecret {
		t.Errorf("expected hashes to match after exclusion, got %q vs %q", h, hNoSecret)
	}
}

func TestCompute_IncludeKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	opts := envhash.Options{IncludeKeys: []string{"A", "C"}}

	h, err := envhash.Compute(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	envSubset := map[string]string{"A": "1", "C": "3"}
	hSubset, _ := envhash.Compute(envSubset, envhash.DefaultOptions())

	if h != hSubset {
		t.Errorf("expected hashes to match for included subset, got %q vs %q", h, hSubset)
	}
}

func TestEqual_SameEnv(t *testing.T) {
	env := map[string]string{"X": "hello", "Y": "world"}
	eq, err := envhash.Equal(env, env, envhash.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eq {
		t.Error("expected Equal to return true for identical maps")
	}
}

func TestEqual_DifferentEnv(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	eq, err := envhash.Equal(a, b, envhash.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eq {
		t.Error("expected Equal to return false for different maps")
	}
}

func TestCompute_EmptyMap(t *testing.T) {
	h, err := envhash.Compute(map[string]string{}, envhash.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == "" {
		t.Error("expected non-empty hash for empty map")
	}
}
