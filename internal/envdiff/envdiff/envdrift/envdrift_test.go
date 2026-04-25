package envdrift

import (
	"strings"
	"testing"
)

func TestDetect_NoChanges(t *testing.T) {
	snap := map[string]string{"A": "1", "B": "2"}
	live := map[string]string{"A": "1", "B": "2"}
	entries := Detect(snap, live, DefaultOptions())
	if len(entries) != 0 {
		t.Fatalf("expected no drift, got %d entries", len(entries))
	}
}

func TestDetect_AddedKey(t *testing.T) {
	snap := map[string]string{"A": "1"}
	live := map[string]string{"A": "1", "B": "2"}
	entries := Detect(snap, live, DefaultOptions())
	if len(entries) != 1 || entries[0].Status != StatusAdded || entries[0].Key != "B" {
		t.Fatalf("expected one added entry for B, got %+v", entries)
	}
}

func TestDetect_RemovedKey(t *testing.T) {
	snap := map[string]string{"A": "1", "B": "2"}
	live := map[string]string{"A": "1"}
	entries := Detect(snap, live, DefaultOptions())
	if len(entries) != 1 || entries[0].Status != StatusRemoved || entries[0].Key != "B" {
		t.Fatalf("expected one removed entry for B, got %+v", entries)
	}
	if entries[0].SnapshotValue != "2" {
		t.Errorf("expected snapshot value '2', got %q", entries[0].SnapshotValue)
	}
}

func TestDetect_ChangedKey(t *testing.T) {
	snap := map[string]string{"A": "old"}
	live := map[string]string{"A": "new"}
	entries := Detect(snap, live, DefaultOptions())
	if len(entries) != 1 || entries[0].Status != StatusChanged {
		t.Fatalf("expected one changed entry, got %+v", entries)
	}
	if entries[0].SnapshotValue != "old" || entries[0].LiveValue != "new" {
		t.Errorf("unexpected values: %+v", entries[0])
	}
}

func TestDetect_IncludeUnchanged(t *testing.T) {
	snap := map[string]string{"A": "1"}
	live := map[string]string{"A": "1"}
	opts := Options{IncludeUnchanged: true}
	entries := Detect(snap, live, opts)
	if len(entries) != 1 || entries[0].Status != StatusUnchanged {
		t.Fatalf("expected one unchanged entry, got %+v", entries)
	}
}

func TestHasDrift_False(t *testing.T) {
	if HasDrift(nil) {
		t.Error("expected no drift for nil entries")
	}
}

func TestHasDrift_True(t *testing.T) {
	entries := []Entry{{Key: "X", Status: StatusAdded, LiveValue: "v"}}
	if !HasDrift(entries) {
		t.Error("expected drift to be detected")
	}
}

func TestFormatText_EmptyEntries(t *testing.T) {
	out := FormatText(nil)
	if !strings.Contains(out, "no drift") {
		t.Errorf("expected 'no drift' message, got %q", out)
	}
}

func TestFormatText_ContainsSymbols(t *testing.T) {
	entries := []Entry{
		{Key: "A", Status: StatusAdded, LiveValue: "1"},
		{Key: "B", Status: StatusRemoved, SnapshotValue: "2"},
		{Key: "C", Status: StatusChanged, SnapshotValue: "old", LiveValue: "new"},
	}
	out := FormatText(entries)
	if !strings.Contains(out, "+ A") {
		t.Errorf("missing added line: %q", out)
	}
	if !strings.Contains(out, "- B") {
		t.Errorf("missing removed line: %q", out)
	}
	if !strings.Contains(out, "~ C") {
		t.Errorf("missing changed line: %q", out)
	}
}
