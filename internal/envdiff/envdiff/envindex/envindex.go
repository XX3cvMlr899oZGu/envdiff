// Package envindex builds an inverted index from a set of env maps,
// mapping each value to the keys that hold it across environments.
package envindex

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultOptions returns a zeroed Options ready for use.
func DefaultOptions() Options {
	return Options{}
}

// Options controls how the index is built.
type Options struct {
	// CaseFold normalises values to lowercase before indexing.
	CaseFold bool
	// ExcludeUnique omits values that appear in only one environment.
	ExcludeUnique bool
}

// Entry records which environments share a particular value.
type Entry struct {
	Value string
	// Occurrences maps env-label → sorted list of keys that hold Value.
	Occurrences map[string][]string
}

// Build constructs an inverted index from the supplied envs map.
// envs maps an environment label (e.g. "prod", "staging") to its key/value pairs.
func Build(envs map[string]map[string]string, opts Options) []Entry {
	// value → envLabel → []key
	idx := map[string]map[string][]string{}

	for label, env := range envs {
		for k, v := range env {
			if opts.CaseFold {
				v = strings.ToLower(v)
			}
			if idx[v] == nil {
				idx[v] = map[string][]string{}
			}
			idx[v][label] = append(idx[v][label], k)
		}
	}

	var entries []Entry
	for val, occ := range idx {
		if opts.ExcludeUnique && len(occ) < 2 {
			continue
		}
		// sort keys within each env for determinism
		for lbl := range occ {
			sort.Strings(occ[lbl])
		}
		entries = append(entries, Entry{Value: val, Occurrences: occ})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value < entries[j].Value
	})
	return entries
}

// FormatText returns a human-readable representation of the index.
func FormatText(entries []Entry) string {
	if len(entries) == 0 {
		return "(empty index)\n"
	}
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("value=%q\n", e.Value))
		labels := make([]string, 0, len(e.Occurrences))
		for lbl := range e.Occurrences {
			labels = append(labels, lbl)
		}
		sort.Strings(labels)
		for _, lbl := range labels {
			sb.WriteString(fmt.Sprintf("  [%s] %s\n", lbl, strings.Join(e.Occurrences[lbl], ", ")))
		}
	}
	return sb.String()
}
