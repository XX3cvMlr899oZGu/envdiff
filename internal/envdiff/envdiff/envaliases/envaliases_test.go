package envaliases_test

import (
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envaliases"
)

func TestApply_NoAliases(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := envaliases.Apply(env, envaliases.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged map, got %v", out)
	}
}

func TestApply_RenamesKey(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value1", "KEEP": "value2"}
	opts := envaliases.Options{
		Aliases:      map[string]string{"OLD_KEY": "NEW_KEY"},
		KeepOriginal: false,
	}
	out, err := envaliases.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", out["NEW_KEY"])
	}
	if out["KEEP"] != "value2" {
		t.Errorf("expected KEEP=value2, got %q", out["KEEP"])
	}
}

func TestApply_KeepOriginal(t *testing.T) {
	env := map[string]string{"SRC": "hello"}
	opts := envaliases.Options{
		Aliases:      map[string]string{"SRC": "DST"},
		KeepOriginal: true,
	}
	out, err := envaliases.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SRC"] != "hello" {
		t.Errorf("expected SRC preserved, got %q", out["SRC"])
	}
	if out["DST"] != "hello" {
		t.Errorf("expected DST=hello, got %q", out["DST"])
	}
}

func TestApply_ConflictReturnsError(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	opts := envaliases.Options{
		Aliases:      map[string]string{"A": "B"},
		KeepOriginal: false,
	}
	_, err := envaliases.Apply(env, opts)
	if err == nil {
		t.Error("expected conflict error, got nil")
	}
}

func TestFormatChanges_ListsSubstitutions(t *testing.T) {
	original := map[string]string{"OLD": "v"}
	aliased := map[string]string{"NEW": "v"}
	aliases := map[string]string{"OLD": "NEW"}
	lines := envaliases.FormatChanges(original, aliased, aliases)
	if len(lines) != 1 {
		t.Fatalf("expected 1 change line, got %d", len(lines))
	}
	if lines[0] != "OLD -> NEW" {
		t.Errorf("unexpected line: %q", lines[0])
	}
}
