package envtrim_test

import (
	"testing"

	"github.com/user/envdiff/internal/envtrim"
)

func TestApply_TrimValues(t *testing.T) {
	env := map[string]string{
		"KEY": "  hello  ",
		"OTHER": "\tworld\n",
	}
	opts := envtrim.DefaultOptions()
	out := envtrim.Apply(env, opts)
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected 'world', got %q", out["OTHER"])
	}
}

func TestApply_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "App_Port": "8080"}
	opts := envtrim.Options{TrimKeys: true, UppercaseKeys: true, TrimValues: false}
	out := envtrim.Apply(env, opts)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key")
	}
	if _, ok := out["APP_PORT"]; !ok {
		t.Error("expected APP_PORT key")
	}
}

func TestApply_TrimKeys(t *testing.T) {
	env := map[string]string{" KEY ": "value"}
	opts := envtrim.Options{TrimKeys: true}
	out := envtrim.Apply(env, opts)
	if _, ok := out["KEY"]; !ok {
		t.Error("expected trimmed key 'KEY'")
	}
}

func TestApply_EmptyKeyDropped(t *testing.T) {
	env := map[string]string{"  ": "value"}
	opts := envtrim.Options{TrimKeys: true}
	out := envtrim.Apply(env, opts)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"KEY": "  val  "}
	opts := envtrim.DefaultOptions()
	envtrim.Apply(env, opts)
	if env["KEY"] != "  val  " {
		t.Error("original map was mutated")
	}
}

func TestApply_NoOptions(t *testing.T) {
	env := map[string]string{"KEY": "  val  "}
	out := envtrim.Apply(env, envtrim.Options{})
	if out["KEY"] != "  val  " {
		t.Errorf("expected untouched value, got %q", out["KEY"])
	}
}
