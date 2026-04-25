// Package envaliases provides key aliasing for env maps.
// It allows renaming keys using a declared alias map while
// preserving original values and detecting conflicts.
package envaliases

import "fmt"

// Options configures alias resolution behaviour.
type Options struct {
	// Aliases maps original key names to their alias names.
	Aliases map[string]string
	// KeepOriginal retains the original key alongside the alias.
	KeepOriginal bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Aliases:      make(map[string]string),
		KeepOriginal: false,
	}
}

// Apply returns a new map with keys renamed according to opts.Aliases.
// If two aliased keys would collide in the output, an error is returned.
func Apply(env map[string]string, opts Options) (map[string]string, error) {
	if len(opts.Aliases) == 0 {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[k] = v
		}
		return out, nil
	}

	out := make(map[string]string, len(env))

	for k, v := range env {
		alias, hasAlias := opts.Aliases[k]
		if hasAlias {
			if _, exists := out[alias]; exists {
				return nil, fmt.Errorf("envaliases: alias %q conflicts with existing key", alias)
			}
			out[alias] = v
			if opts.KeepOriginal {
				out[k] = v
			}
		} else {
			if _, exists := out[k]; exists {
				return nil, fmt.Errorf("envaliases: key %q already present in output", k)
			}
			out[k] = v
		}
	}

	return out, nil
}

// FormatChanges returns a human-readable list of alias substitutions applied.
func FormatChanges(original, aliased map[string]string, aliases map[string]string) []string {
	var lines []string
	for orig, alias := range aliases {
		if _, ok := original[orig]; ok {
			lines = append(lines, fmt.Sprintf("%s -> %s", orig, alias))
		}
	}
	return lines
}
