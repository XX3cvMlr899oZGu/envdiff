// Package envdrift detects configuration drift between a saved snapshot
// and a live environment map, categorising each key by its drift status.
package envdrift

import (
	"fmt"
	"sort"
	"strings"
)

// Status describes how a key has drifted from the snapshot.
type Status string

const (
	StatusAdded    Status = "added"    // present in live, absent in snapshot
	StatusRemoved  Status = "removed"  // present in snapshot, absent in live
	StatusChanged  Status = "changed"  // present in both, values differ
	StatusUnchanged Status = "unchanged" // present in both, values equal
)

// Entry holds drift information for a single key.
type Entry struct {
	Key           string
	Status        Status
	SnapshotValue string
	LiveValue     string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{IncludeUnchanged: false}
}

// Options controls Detect behaviour.
type Options struct {
	IncludeUnchanged bool
}

// Detect compares snapshot against live and returns drift entries.
func Detect(snapshot, live map[string]string, opts Options) []Entry {
	seen := make(map[string]bool)
	var entries []Entry

	for k, sv := range snapshot {
		seen[k] = true
		if lv, ok := live[k]; !ok {
			entries = append(entries, Entry{Key: k, Status: StatusRemoved, SnapshotValue: sv})
		} else if lv != sv {
			entries = append(entries, Entry{Key: k, Status: StatusChanged, SnapshotValue: sv, LiveValue: lv})
		} else if opts.IncludeUnchanged {
			entries = append(entries, Entry{Key: k, Status: StatusUnchanged, SnapshotValue: sv, LiveValue: lv})
		}
	}

	for k, lv := range live {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: StatusAdded, LiveValue: lv})
		}
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	return entries
}

// HasDrift returns true if any entry represents an actual change.
func HasDrift(entries []Entry) bool {
	for _, e := range entries {
		if e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}

// FormatText renders drift entries as a human-readable string.
func FormatText(entries []Entry) string {
	if len(entries) == 0 {
		return "no drift detected\n"
	}
	var sb strings.Builder
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			fmt.Fprintf(&sb, "+ %s = %q\n", e.Key, e.LiveValue)
		case StatusRemoved:
			fmt.Fprintf(&sb, "- %s (was %q)\n", e.Key, e.SnapshotValue)
		case StatusChanged:
			fmt.Fprintf(&sb, "~ %s: %q → %q\n", e.Key, e.SnapshotValue, e.LiveValue)
		case StatusUnchanged:
			fmt.Fprintf(&sb, "  %s = %q\n", e.Key, e.LiveValue)
		}
	}
	return sb.String()
}
