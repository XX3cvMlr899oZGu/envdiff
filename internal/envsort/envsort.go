// Package envsort provides utilities for sorting env maps into ordered slices.
package envsort

import (
	"sort"
	"strings"
)

// SortOrder defines how keys should be sorted.
type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

// Options configures the sort behaviour.
type Options struct {
	Order  SortOrder
	Prefix string // if set, keys with this prefix are sorted first
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Order: Ascending}
}

// Entry is a key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Apply sorts the env map into a slice of Entry according to opts.
func Apply(env map[string]string, opts Options) []Entry {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		ai := opts.Prefix != "" && strings.HasPrefix(keys[i], opts.Prefix)
		aj := opts.Prefix != "" && strings.HasPrefix(keys[j], opts.Prefix)
		if ai != aj {
			return ai // prefix keys come first
		}
		if opts.Order == Descending {
			return keys[i] > keys[j]
		}
		return keys[i] < keys[j]
	})

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Value: env[k]})
	}
	return entries
}
