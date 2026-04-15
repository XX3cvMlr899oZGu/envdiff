// Package rename provides utilities for renaming keys in an env map,
// either via explicit mappings or a prefix substitution.
package rename

import "fmt"

// Options controls how keys are renamed.
type Options struct {
	// Map is an explicit old->new key mapping.
	Map map[string]string
	// OldPrefix and NewPrefix, when both non-empty, replace a key prefix.
	OldPrefix string
	NewPrefix string
}

// DefaultOptions returns an Options with an empty map and no prefix substitution.
func DefaultOptions() Options {
	return Options{Map: make(map[string]string)}
}

// Result holds the renamed env map and a log of changes made.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Change records a single key rename.
type Change struct {
	OldKey string
	NewKey string
}

// Apply renames keys in env according to opts and returns a new map.
// Returns an error if a rename would produce a duplicate key.
func Apply(env map[string]string, opts Options) (Result, error) {
	out := make(map[string]string, len(env))
	var changes []Change

	for k, v := range env {
		newKey := k

		if mapped, ok := opts.Map[k]; ok {
			newKey = mapped
		} else if opts.OldPrefix != "" && opts.NewPrefix != "" {
			if len(k) >= len(opts.OldPrefix) && k[:len(opts.OldPrefix)] == opts.OldPrefix {
				newKey = opts.NewPrefix + k[len(opts.OldPrefix):]
			}
		}

		if _, exists := out[newKey]; exists {
			return Result{}, fmt.Errorf("rename conflict: key %q already exists in output", newKey)
		}

		out[newKey] = v
		if newKey != k {
			changes = append(changes, Change{OldKey: k, NewKey: newKey})
		}
	}

	return Result{Env: out, Changes: changes}, nil
}

// FormatChanges returns a human-readable summary of all renames.
func FormatChanges(changes []Change) string {
	if len(changes) == 0 {
		return "no keys renamed\n"
	}
	out := ""
	for _, c := range changes {
		out += fmt.Sprintf("  %s -> %s\n", c.OldKey, c.NewKey)
	}
	return out
}
