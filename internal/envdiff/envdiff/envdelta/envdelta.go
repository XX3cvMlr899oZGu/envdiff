// Package envdelta computes a structured delta between two env maps,
// categorising every key as added, removed, changed, or unchanged.
package envdelta

import "sort"

// Status represents the kind of change for a single key.
type Status string

const (
	StatusAdded     Status = "added"
	StatusRemoved   Status = "removed"
	StatusChanged   Status = "changed"
	StatusUnchanged Status = "unchanged"
)

// Entry holds the delta information for one key.
type Entry struct {
	Key    string
	Status Status
	OldVal string
	NewVal string
}

// Delta is the full result of comparing two env maps.
type Delta struct {
	Entries []Entry
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options { return Options{IncludeUnchanged: false} }

// Options controls what Compute includes in the result.
type Options struct {
	IncludeUnchanged bool
}

// Compute returns a Delta describing every key difference between base and next.
func Compute(base, next map[string]string, opts Options) Delta {
	seen := make(map[string]bool)
	var entries []Entry

	for k, oldVal := range base {
		seen[k] = true
		if newVal, ok := next[k]; ok {
			if oldVal == newVal {
				if opts.IncludeUnchanged {
					entries = append(entries, Entry{Key: k, Status: StatusUnchanged, OldVal: oldVal, NewVal: newVal})
				}
			} else {
				entries = append(entries, Entry{Key: k, Status: StatusChanged, OldVal: oldVal, NewVal: newVal})
			}
		} else {
			entries = append(entries, Entry{Key: k, Status: StatusRemoved, OldVal: oldVal})
		}
	}

	for k, newVal := range next {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: StatusAdded, NewVal: newVal})
		}
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	return Delta{Entries: entries}
}

// HasChanges reports whether the delta contains any added, removed, or changed entries.
func (d Delta) HasChanges() bool {
	for _, e := range d.Entries {
		if e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}

// ByStatus returns only the entries that match the given status.
func (d Delta) ByStatus(s Status) []Entry {
	var out []Entry
	for _, e := range d.Entries {
		if e.Status == s {
			out = append(out, e)
		}
	}
	return out
}
