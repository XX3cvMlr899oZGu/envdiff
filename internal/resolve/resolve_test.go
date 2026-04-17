package resolve_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/resolve"
)

func TestResolve_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	got, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", got)
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE": "/home/user",
		"PATH": "${BASE}/bin",
	}
	got, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["PATH"] != "/home/user/bin" {
		t.Errorf("expected /home/user/bin, got %q", got["PATH"])
	}
}

func TestResolve_DollarStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://$HOST:8080",
	}
	got, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["URL"] != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %q", got["URL"])
	}
}

func TestResolve_ChainedReferences(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	got, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["C"] != "hello_world!" {
		t.Errorf("expected hello_world!, got %q", got["C"])
	}
}

func TestResolve_UnresolvableLeftAsIs(t *testing.T) {
	env := map[string]string{"FOO": "${MISSING}/path"}
	got, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "${MISSING}/path" {
		t.Errorf("expected reference left as-is, got %q", got["FOO"])
	}
}

func TestResolve_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"A": "x", "B": "${A}y"}
	_, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["B"] != "${A}y" {
		t.Errorf("original map mutated: B = %q", env["B"])
	}
}

func TestResolve_InvalidMaxDepth(t *testing.T) {
	env := map[string]string{"A": "1"}
	_, err := resolve.Resolve(env, resolve.Options{MaxDepth: 0})
	if err == nil {
		t.Error("expected error for MaxDepth=0, got nil")
	}
}

func TestResolve_CircularReference(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	_, err := resolve.Resolve(env, resolve.DefaultOptions())
	if err == nil {
		t.Error("expected error for circular reference, got nil")
	}
}
