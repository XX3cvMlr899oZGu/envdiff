// Package envclone provides utilities for deep-copying env maps
// with optional key transformation and filtering.
package envclone

import "strings"

// Options controls how the clone is performed.
type Options struct {
	// KeyPrefix filters keys by prefix before cloning.
	KeyPrefix string
	// StripPrefix removes the prefix from keys in the output.
	StripPrefix bool
	// KeyTransform applies a function to each key.
	KeyTransform func(string) string
}

// DefaultOptions returns Options with no transformations.
func DefaultOptions() Options {
	return Options{}
}

// Apply clones src into a new map, applying options.
func Apply(src map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if opts.KeyPrefix != "" && !strings.HasPrefix(k, opts.KeyPrefix) {
			continue
		}
		key := k
		if opts.StripPrefix && opts.KeyPrefix != "" {
			key = strings.TrimPrefix(k, opts.KeyPrefix)
		}
		if opts.KeyTransform != nil {
			key = opts.KeyTransform(key)
		}
		out[key] = v
	}
	return out
}
