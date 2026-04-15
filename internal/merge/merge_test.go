package merge

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	res, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if res.Env["FOO"] != "1" || res.Env["BAR"] != "2" || res.Env["BAZ"] != "3" {
		t.Errorf("unexpected merged env: %v", res.Env)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override"}

	res, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "original" {
		t.Errorf("expected 'original', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "KEY" {
		t.Errorf("expected conflict on KEY, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override"}

	res, err := Merge(StrategyLast, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "override" {
		t.Errorf("expected 'override', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyError(t *testing.T) {
	a := map[string]string{"KEY": "a"}
	b := map[string]string{"KEY": "b"}

	_, err := Merge(StrategyError, a, b)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_SameValueNoConflict(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}

	res, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("identical values should not be conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	res, err := Merge(StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}
