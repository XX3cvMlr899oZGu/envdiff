package envdelta

import (
	"testing"
)

func TestCompute_NoChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "2"}
	d := Compute(base, next, DefaultOptions())
	if d.HasChanges() {
		t.Fatal("expected no changes")
	}
	if len(d.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(d.Entries))
	}
}

func TestCompute_AddedKey(t *testing.T) {
	base := map[string]string{"A": "1"}
	next := map[string]string{"A": "1", "B": "2"}
	d := Compute(base, next, DefaultOptions())
	if !d.HasChanges() {
		t.Fatal("expected changes")
	}
	added := d.ByStatus(StatusAdded)
	if len(added) != 1 || added[0].Key != "B" {
		t.Fatalf("unexpected added entries: %v", added)
	}
}

func TestCompute_RemovedKey(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1"}
	d := Compute(base, next, DefaultOptions())
	removed := d.ByStatus(StatusRemoved)
	if len(removed) != 1 || removed[0].Key != "B" || removed[0].OldVal != "2" {
		t.Fatalf("unexpected removed entries: %v", removed)
	}
}

func TestCompute_ChangedKey(t *testing.T) {
	base := map[string]string{"A": "old"}
	next := map[string]string{"A": "new"}
	d := Compute(base, next, DefaultOptions())
	changed := d.ByStatus(StatusChanged)
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(changed))
	}
	if changed[0].OldVal != "old" || changed[0].NewVal != "new" {
		t.Fatalf("wrong values: %+v", changed[0])
	}
}

func TestCompute_IncludeUnchanged(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "2"}
	opts := Options{IncludeUnchanged: true}
	d := Compute(base, next, opts)
	if len(d.Entries) != 2 {
		t.Fatalf("expected 2 unchanged entries, got %d", len(d.Entries))
	}
	for _, e := range d.Entries {
		if e.Status != StatusUnchanged {
			t.Fatalf("expected unchanged, got %s", e.Status)
		}
	}
}

func TestCompute_SortedByKey(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2", "M": "3"}
	next := map[string]string{"Z": "1", "A": "9", "M": "3"}
	opts := Options{IncludeUnchanged: true}
	d := Compute(base, next, opts)
	keys := make([]string, len(d.Entries))
	for i, e := range d.Entries {
		keys[i] = e.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Fatalf("entries not sorted: %v", keys)
		}
	}
}

func TestByStatus_Empty(t *testing.T) {
	d := Delta{}
	if got := d.ByStatus(StatusAdded); len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}
