// Package resolve provides utilities for resolving environment variable
// references within .env file values (e.g. VAR=${OTHER_VAR}).
package resolve

import (
	"fmt"
	"regexp"
	"strings"
)

// refPattern matches ${VAR} and $VAR style references.
var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls resolution behaviour.
type Options struct {
	// MaxDepth limits recursive resolution passes to prevent cycles.
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{MaxDepth: 10}
}

// Resolve expands variable references in all values of env using the same map
// as the source of substitutions. Unresolvable references are left as-is.
// Returns a new map; the original is not mutated.
func Resolve(env map[string]string, opts Options) (map[string]string, error) {
	if opts.MaxDepth <= 0 {
		return nil, fmt.Errorf("resolve: MaxDepth must be greater than zero")
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	for i := 0; i < opts.MaxDepth; i++ {
		changed := false
		for k, v := range result {
			expanded := expand(v, result)
			if expanded != v {
				result[k] = expanded
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	return result, nil
}

// expand replaces all variable references in s using lookup.
func expand(s string, lookup map[string]string) string {
	return refPattern.ReplaceAllStringFunc(s, func(match string) string {
		name := strings.TrimPrefix(strings.TrimPrefix(strings.Trim(match, "${}"), "${"), "$")
		name = strings.TrimSuffix(name, "}")
		if val, ok := lookup[name]; ok {
			return val
		}
		return match
	})
}
