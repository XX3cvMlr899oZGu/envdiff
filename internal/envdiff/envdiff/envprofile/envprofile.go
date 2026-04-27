// Package envprofile selects a named subset of environment variables
// based on a profile definition (e.g. "production", "staging").
package envprofile

import (
	"fmt"
	"sort"
	"strings"
)

// Profile defines a named set of key patterns that belong together.
type Profile struct {
	Name     string
	Patterns []string // exact keys or prefix globs ending with "*"
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{IncludeUnmatched: false}
}

// Options controls Apply behaviour.
type Options struct {
	// IncludeUnmatched keeps keys that do not match any profile pattern.
	IncludeUnmatched bool
}

// Apply filters env to only the keys matched by profile p.
// If opts.IncludeUnmatched is true, non-matching keys are kept as-is.
func Apply(env map[string]string, p Profile, opts Options) (map[string]string, error) {
	if p.Name == "" {
		return nil, fmt.Errorf("envprofile: profile name must not be empty")
	}
	result := make(map[string]string, len(env))
	for k, v := range env {
		if matchesAny(k, p.Patterns) {
			result[k] = v
		} else if opts.IncludeUnmatched {
			result[k] = v
		}
	}
	return result, nil
}

// MatchedKeys returns the sorted list of keys from env that match profile p.
func MatchedKeys(env map[string]string, p Profile) []string {
	var keys []string
	for k := range env {
		if matchesAny(k, p.Patterns) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func matchesAny(key string, patterns []string) bool {
	for _, pat := range patterns {
		if strings.HasSuffix(pat, "*") {
			if strings.HasPrefix(key, strings.TrimSuffix(pat, "*")) {
				return true
			}
		} else if key == pat {
			return true
		}
	}
	return false
}
