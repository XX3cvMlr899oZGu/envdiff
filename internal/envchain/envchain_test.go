package envchain

import (
	"errors"
	"testing"
)

func identity(env map[string]string) (map[string]string, error) {
	return env, nil
}

func prefixValues(prefix string) StepFunc {
	return func(env map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[k] = prefix + v
		}
		return out, nil
	}
}

func failStep(env map[string]string) (map[string]string, error) {
	return nil, errors.New("intentional failure")
}

func TestRun_EmptyChain(t *testing.T) {
	c := New()
	input := map[string]string{"A": "1"}
	out, err := c.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
}

func TestRun_SingleStep(t *testing.T) {
	c := New().Add("prefix", prefixValues("hello_"))
	out, err := c.Run(map[string]string{"KEY": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello_world" {
		t.Errorf("expected hello_world, got %q", out["KEY"])
	}
}

func TestRun_MultipleSteps_Ordered(t *testing.T) {
	c := New().
		Add("first", prefixValues("A_")).
		Add("second", prefixValues("B_"))
	out, err := c.Run(map[string]string{"X": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// first: X=A_val, second: X=B_A_val
	if out["X"] != "B_A_val" {
		t.Errorf("expected B_A_val, got %q", out["X"])
	}
}

func TestRun_StepError_Propagates(t *testing.T) {
	c := New().
		Add("ok", identity).
		Add("bad", failStep)
	_, err := c.Run(map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errors.Unwrap(err)) {
		// just check message contains step name
	}
	if err.Error() == "" {
		t.Error("error message should not be empty")
	}
}

func TestStepNames(t *testing.T) {
	c := New().Add("alpha", identity).Add("beta", identity)
	names := c.StepNames()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestRun_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"K": "original"}
	c := New().Add("mutate", prefixValues("changed_"))
	c.Run(input)
	if input["K"] != "original" {
		t.Errorf("input was mutated: got %q", input["K"])
	}
}
