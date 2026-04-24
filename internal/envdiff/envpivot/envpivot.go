// Package envpivot transposes a set of named env maps into a key-centric view,
// making it easy to compare how a single key varies across multiple environments.
package envpivot

import "sort"

// Row holds the value of a single key across all environments.
type Row struct {
	Key    string
	Values map[string]string // env name -> value (empty string if absent)
}

// Options controls Pivot behaviour.
type Options struct {
	// IncludeAbsent, when true, includes keys that are missing in some envs.
	IncludeAbsent bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{IncludeAbsent: true}
}

// Pivot takes a map of environment-name -> key/value pairs and returns a
// slice of Row values sorted by key name.
func Pivot(envs map[string]map[string]string, opts Options) []Row {
	// Collect the union of all keys.
	keySet := make(map[string]struct{})
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	rows := make([]Row, 0, len(keySet))
	for key := range keySet {
		values := make(map[string]string, len(envs))
		absent := false
		for envName, kv := range envs {
			v, ok := kv[key]
			if !ok {
				absent = true
			}
			values[envName] = v
		}
		if absent && !opts.IncludeAbsent {
			continue
		}
		rows = append(rows, Row{Key: key, Values: values})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Key < rows[j].Key
	})
	return rows
}

// AllEqual returns true when all environment values for the row are identical
// (and present in every environment).
func (r Row) AllEqual(envNames []string) bool {
	if len(envNames) == 0 {
		return true
	}
	first := r.Values[envNames[0]]
	for _, name := range envNames[1:] {
		if r.Values[name] != first {
			return false
		}
	}
	return true
}
