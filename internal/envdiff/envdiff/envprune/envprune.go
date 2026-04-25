// Package envprune removes keys from an env map based on configurable criteria
// such as empty values, duplicate values, or matching patterns.
package envprune

import (
	"fmt"
	"regexp"
)

// DefaultOptions returns an Options with safe defaults.
func DefaultOptions() Options {
	return Options{
		RemoveEmpty:     false,
		RemoveDuplicates: false,
	}
}

// Options controls which keys are pruned.
type Options struct {
	// RemoveEmpty drops keys whose value is the empty string.
	RemoveEmpty bool

	// RemoveDuplicates drops keys whose value is identical to a previously seen
	// value (first occurrence is kept).
	RemoveDuplicates bool

	// PatternKeys is an optional regex; keys matching it are dropped.
	PatternKeys string
}

// Result holds the pruned map and the list of keys that were removed.
type Result struct {
	Env     map[string]string
	Removed []string
}

// Apply prunes env according to opts and returns a Result.
// The original map is never mutated.
func Apply(env map[string]string, opts Options) (Result, error) {
	var keyRe *regexp.Regexp
	if opts.PatternKeys != "" {
		var err error
		keyRe, err = regexp.Compile(opts.PatternKeys)
		if err != nil {
			return Result{}, fmt.Errorf("envprune: invalid pattern_keys regexp: %w", err)
		}
	}

	seen := map[string]bool{} // value → already present
	out := make(map[string]string, len(env))
	var removed []string

	for k, v := range env {
		switch {
		case opts.RemoveEmpty && v == "":
			removed = append(removed, k)
		case opts.RemoveDuplicates && seen[v]:
			removed = append(removed, k)
		case keyRe != nil && keyRe.MatchString(k):
			removed = append(removed, k)
		default:
			out[k] = v
			seen[v] = true
		}
	}

	return Result{Env: out, Removed: removed}, nil
}
