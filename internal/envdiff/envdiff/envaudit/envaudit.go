// Package envaudit provides audit trail functionality for tracking changes
// to environment variable sets over time, recording who changed what and when.
package envaudit

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ChangeKind describes the type of change recorded in an audit entry.
type ChangeKind string

const (
	ChangeAdded   ChangeKind = "added"
	ChangeRemoved ChangeKind = "removed"
	ChangeUpdated ChangeKind = "updated"
)

// Entry represents a single audited change to an environment variable.
type Entry struct {
	Timestamp time.Time  `json:"timestamp"`
	Key       string     `json:"key"`
	Kind      ChangeKind `json:"kind"`
	OldValue  string     `json:"old_value,omitempty"`
	NewValue  string     `json:"new_value,omitempty"`
	Author    string     `json:"author,omitempty"`
	Note      string     `json:"note,omitempty"`
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		RedactValues: true,
		RedactPlaceholder: "[REDACTED]",
	}
}

// Options controls audit behaviour.
type Options struct {
	// RedactValues replaces sensitive values with a placeholder in audit entries.
	RedactValues bool
	// RedactPlaceholder is the string used when RedactValues is true.
	RedactPlaceholder string
	// Author is an optional label (e.g. username or service name) attached to each entry.
	Author string
	// Note is an optional free-text annotation attached to each entry.
	Note string
}

// Compute produces an audit trail by comparing a previous env map (before) with
// a current env map (after). Each key that was added, removed, or updated
// produces one Entry. The entries are sorted deterministically by key.
func Compute(before, after map[string]string, opts Options) []Entry {
	if opts.RedactPlaceholder == "" {
		opts.RedactPlaceholder = "[REDACTED]"
	}

	now := time.Now().UTC()
	var entries []Entry

	// Detect removed and updated keys.
	for k, oldVal := range before {
		newVal, exists := after[k]
		if !exists {
			entries = append(entries, Entry{
				Timestamp: now,
				Key:       k,
				Kind:      ChangeRemoved,
				OldValue:  redact(oldVal, opts),
				Author:    opts.Author,
				Note:      opts.Note,
			})
		} else if newVal != oldVal {
			entries = append(entries, Entry{
				Timestamp: now,
				Key:       k,
				Kind:      ChangeUpdated,
				OldValue:  redact(oldVal, opts),
				NewValue:  redact(newVal, opts),
				Author:    opts.Author,
				Note:      opts.Note,
			})
		}
	}

	// Detect added keys.
	for k, newVal := range after {
		if _, exists := before[k]; !exists {
			entries = append(entries, Entry{
				Timestamp: now,
				Key:       k,
				Kind:      ChangeAdded,
				NewValue:  redact(newVal, opts),
				Author:    opts.Author,
				Note:      opts.Note,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Key != entries[j].Key {
			return entries[i].Key < entries[j].Key
		}
		return string(entries[i].Kind) < string(entries[j].Kind)
	})

	return entries
}

// HasChanges returns true when the audit trail contains at least one entry.
func HasChanges(entries []Entry) bool {
	return len(entries) > 0
}

// FormatText returns a human-readable summary of the audit trail.
func FormatText(entries []Entry) string {
	if len(entries) == 0 {
		return "no changes detected"
	}

	var sb strings.Builder
	for _, e := range entries {
		ts := e.Timestamp.Format(time.RFC3339)
		switch e.Kind {
		case ChangeAdded:
			fmt.Fprintf(&sb, "[%s] + %s = %q", ts, e.Key, e.NewValue)
		case ChangeRemoved:
			fmt.Fprintf(&sb, "[%s] - %s (was %q)", ts, e.Key, e.OldValue)
		case ChangeUpdated:
			fmt.Fprintf(&sb, "[%s] ~ %s: %q → %q", ts, e.Key, e.OldValue, e.NewValue)
		}
		if e.Author != "" {
			fmt.Fprintf(&sb, " by %s", e.Author)
		}
		if e.Note != "" {
			fmt.Fprintf(&sb, " (%s)", e.Note)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// redact returns the placeholder when RedactValues is enabled, otherwise the
// original value is returned unchanged.
func redact(value string, opts Options) string {
	if opts.RedactValues {
		return opts.RedactPlaceholder
	}
	return value
}
