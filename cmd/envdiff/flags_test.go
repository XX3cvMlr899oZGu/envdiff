package main

import (
	"flag"
	"os"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestExcludeList_MultipleValues(t *testing.T) {
	var list excludeList
	if err := list.Set("KEY_A"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := list.Set("KEY_B"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
	if list[0] != "KEY_A" || list[1] != "KEY_B" {
		t.Errorf("unexpected values: %v", list)
	}
}

func TestExcludeList_String(t *testing.T) {
	list := excludeList{"KEY_A", "KEY_B"}
	if list.String() != "KEY_A,KEY_B" {
		t.Errorf("unexpected string: %s", list.String())
	}
}

func TestExcludeList_Empty(t *testing.T) {
	var list excludeList
	if list.String() != "" {
		t.Errorf("expected empty string, got %q", list.String())
	}
}
