package cast_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/cast"
)

func TestApply_FallbackToString(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp"}
	opts := cast.DefaultOptions()
	results, err := cast.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].TypeName != "string" || results[0].Value != "myapp" {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

func TestApply_BoolHint(t *testing.T) {
	env := map[string]string{"DEBUG": "true"}
	opts := cast.DefaultOptions()
	opts.TypeHints["DEBUG"] = "bool"
	results, err := cast.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != true {
		t.Errorf("expected true, got %v", results[0].Value)
	}
}

func TestApply_IntHint(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	opts := cast.DefaultOptions()
	opts.TypeHints["PORT"] = "int"
	results, err := cast.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != 8080 {
		t.Errorf("expected 8080, got %v", results[0].Value)
	}
}

func TestApply_FloatHint(t *testing.T) {
	env := map[string]string{"RATIO": "3.14"}
	opts := cast.DefaultOptions()
	opts.TypeHints["RATIO"] = "float"
	results, err := cast.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value.(float64) < 3.13 {
		t.Errorf("unexpected float value: %v", results[0].Value)
	}
}

func TestApply_InvalidBool(t *testing.T) {
	env := map[string]string{"DEBUG": "notabool"}
	opts := cast.DefaultOptions()
	opts.TypeHints["DEBUG"] = "bool"
	_, err := cast.Apply(env, opts)
	if err == nil {
		t.Error("expected error for invalid bool")
	}
}

func TestApply_UnknownTypeHint(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	opts := cast.DefaultOptions()
	opts.TypeHints["KEY"] = "uuid"
	_, err := cast.Apply(env, opts)
	if err == nil {
		t.Error("expected error for unknown type hint")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	results, err := cast.Apply(map[string]string{}, cast.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
