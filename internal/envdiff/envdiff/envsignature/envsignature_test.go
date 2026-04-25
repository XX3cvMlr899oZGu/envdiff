package envsignature_test

import (
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envsignature"
)

func TestCompute_DeterministicOutput(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := envsignature.DefaultOptions()

	r1 := envsignature.Compute(env, opts)
	r2 := envsignature.Compute(env, opts)

	if r1.Signature != r2.Signature {
		t.Errorf("expected deterministic signature, got %q and %q", r1.Signature, r2.Signature)
	}
	if r1.KeyCount != 2 {
		t.Errorf("expected KeyCount=2, got %d", r1.KeyCount)
	}
}

func TestCompute_DifferentEnvsProduceDifferentSignatures(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "baz"}
	opts := envsignature.DefaultOptions()

	if envsignature.Compute(a, opts).Signature == envsignature.Compute(b, opts).Signature {
		t.Error("expected different signatures for different values")
	}
}

func TestCompute_KeysOnlyIgnoresValues(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "1"}
	b := map[string]string{"FOO": "zzz", "BAZ": "999"}
	opts := envsignature.Options{IncludeValues: false}

	if envsignature.Compute(a, opts).Signature != envsignature.Compute(b, opts).Signature {
		t.Error("expected same signature when values differ but keys-only mode is on")
	}
}

func TestCompute_PrefixFilter(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	opts := envsignature.Options{IncludeValues: true, Prefix: "APP_"}

	r := envsignature.Compute(env, opts)
	if r.KeyCount != 2 {
		t.Errorf("expected KeyCount=2 with prefix filter, got %d", r.KeyCount)
	}
}

func TestEqual_SameEnv(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	opts := envsignature.DefaultOptions()

	if !envsignature.Equal(env, env, opts) {
		t.Error("expected Equal to return true for same env")
	}
}

func TestEqual_DifferentEnv(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	opts := envsignature.DefaultOptions()

	if envsignature.Equal(a, b, opts) {
		t.Error("expected Equal to return false for different envs")
	}
}

func TestCompute_EmptyMap(t *testing.T) {
	r := envsignature.Compute(map[string]string{}, envsignature.DefaultOptions())
	if r.KeyCount != 0 {
		t.Errorf("expected KeyCount=0, got %d", r.KeyCount)
	}
	if r.Signature == "" {
		t.Error("expected non-empty signature even for empty map")
	}
}
