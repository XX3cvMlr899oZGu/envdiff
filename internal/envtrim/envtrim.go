// Package envtrim removes leading/trailing whitespace from env values
// and optionally normalizes keys to uppercase.
package envtrim

import (
	"strings"
)

// Options controls trimming behaviour.
type Options struct {
	// TrimValues removes surrounding whitespace from values.
	TrimValues bool
	// UppercaseKeys normalizes all keys to uppercase.
	UppercaseKeys bool
	// TrimKeys removes surrounding whitespace from keys.
	TrimKeys bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		TrimValues: true,
		UppercaseKeys: false,
		TrimKeys: true,
	}
}

// Apply returns a new map with trimming rules applied.
// The original map is never mutated.
func Apply(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.TrimKeys {
			k = strings.TrimSpace(k)
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if k != "" {
			out[k] = v
		}
	}
	return out
}
