package envbenchmark_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envbenchmark"
)

func makeEnv(n int) map[string]string {
	m := make(map[string]string, n)
	for i := 0; i < n; i++ {
		m[strings.Repeat("K", i+1)] = strings.Repeat("v", i+1)
	}
	return m
}

func noop(env map[string]string) error { return nil }

func TestRun_ReturnsResult(t *testing.T) {
	env := makeEnv(5)
	r, err := envbenchmark.Run("noop", env, noop, envbenchmark.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Label != "noop" {
		t.Errorf("expected label 'noop', got %q", r.Label)
	}
	if r.KeyCount != 5 {
		t.Errorf("expected key count 5, got %d", r.KeyCount)
	}
	if r.Duration < 0 {
		t.Error("duration should be non-negative")
	}
}

func TestRun_MultipleRuns(t *testing.T) {
	env := makeEnv(3)
	opts := envbenchmark.Options{Runs: 10}
	r, err := envbenchmark.Run("multi", env, noop, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.KeyCount != 3 {
		t.Errorf("expected 3 keys, got %d", r.KeyCount)
	}
}

func TestRun_FnError_Propagates(t *testing.T) {
	env := makeEnv(2)
	failFn := func(map[string]string) error { return errors.New("boom") }
	_, err := envbenchmark.Run("fail", env, failFn, envbenchmark.DefaultOptions())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("expected 'boom' in error, got %v", err)
	}
}

func TestFormatText_ContainsLabel(t *testing.T) {
	results := []envbenchmark.Result{
		{Label: "alpha", KeyCount: 10},
		{Label: "beta", KeyCount: 5},
	}
	out := envbenchmark.FormatText(results)
	if !strings.Contains(out, "alpha") {
		t.Error("expected 'alpha' in output")
	}
	if !strings.Contains(out, "beta") {
		t.Error("expected 'beta' in output")
	}
}

func TestFormatText_Empty(t *testing.T) {
	out := envbenchmark.FormatText(nil)
	if !strings.Contains(out, "no benchmark results") {
		t.Errorf("unexpected output: %q", out)
	}
}
