// Package envpivot pivots multiple env maps into a table keyed by variable name.
package envpivot

import "sort"

// Row represents a single key across all environments.
type Row struct {
	Key    string
	Values map[string]string // env label -> value
}

// AllEqual returns true if all present values are identical.
func (r Row) AllEqual() bool {
	var seen string
	first := true
	for _, v := range r.Values {
		if first {
			seen = v
			first = false
			continue
		}
		if v != seen {
			return false
		}
	}
	return true
}

// Options controls pivot behaviour.
type Options struct {
	// ExcludeMissing drops rows where any environment is missing the key.
	ExcludeMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{ExcludeMissing: false}
}

// Pivot converts a map of label->env into a slice of Rows sorted by key.
// envs maps an environment label (e.g. "staging") to its parsed key/value map.
func Pivot(envs map[string]map[string]string, opts Options) []Row {
	// collect all keys
	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	var rows []Row
	for key := range keySet {
		row := Row{Key: key, Values: make(map[string]string, len(envs))}
		missingAny := false
		for label, env := range envs {
			v, ok := env[key]
			if !ok {
				missingAny = true
			} else {
				row.Values[label] = v
			}
		}
		if opts.ExcludeMissing && missingAny {
			continue
		}
		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Key < rows[j].Key
	})
	return rows
}
