package envdefault_test

import (
	"sort"
	"testing"

	"github.com/user/envdiff/internal/envdefault"
)

func TestApply_FillsMissingKeys(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"A": "99", "B": "2"}
	out := envdefault.Apply(env, defs, envdefault.DefaultOptions())
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %s", out["B"])
	}
}

func TestApply_Overwrite(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"A": "99"}
	opts := envdefault.Options{Overwrite: true}
	out := envdefault.Apply(env, defs, opts)
	if out["A"] != "99" {
		t.Errorf("expected A=99, got %s", out["A"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"B": "2"}
	envdefault.Apply(env, defs, envdefault.DefaultOptions())
	if _, ok := env["B"]; ok {
		t.Error("original env was mutated")
	}
}

func TestApply_EmptyDefaults(t *testing.T) {
	env := map[string]string{"A": "1"}
	out := envdefault.Apply(env, map[string]string{}, envdefault.DefaultOptions())
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestMissingKeys(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"A": "1", "B": "2", "C": "3"}
	missing := envdefault.MissingKeys(env, defs)
	sort.Strings(missing)
	if len(missing) != 2 || missing[0] != "B" || missing[1] != "C" {
		t.Errorf("unexpected missing keys: %v", missing)
	}
}

func TestMissingKeys_NoneAbsent(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	defs := map[string]string{"A": "x", "B": "y"}
	if len(envdefault.MissingKeys(env, defs)) != 0 {
		t.Error("expected no missing keys")
	}
}
