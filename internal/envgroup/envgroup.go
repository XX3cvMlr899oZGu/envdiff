// Package envgroup groups env keys by prefix into named sub-maps.
package envgroup

import (
	"sort"
	"strings"
)

// Group represents a named collection of env key-value pairs.
type Group struct {
	Name string
	Keys map[string]string
}

// DefaultSeparator is the delimiter used to split prefix from key.
const DefaultSeparator = "_"

// Options controls grouping behaviour.
type Options struct {
	Separator string
	StripPrefix bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator:   DefaultSeparator,
		StripPrefix: false,
	}
}

// Apply partitions env into groups keyed by the first prefix segment.
// Keys with no separator are placed in a group named "".
func Apply(env map[string]string, opts Options) []Group {
	if opts.Separator == "" {
		opts.Separator = DefaultSeparator
	}

	buckets := make(map[string]map[string]string)

	for k, v := range env {
		parts := strings.SplitN(k, opts.Separator, 2)
		prefix := ""
		key := k
		if len(parts) == 2 {
			prefix = parts[0]
			if opts.StripPrefix {
				key = parts[1]
			}
		}
		if buckets[prefix] == nil {
			buckets[prefix] = make(map[string]string)
		}
		buckets[prefix][key] = v
	}

	names := make([]string, 0, len(buckets))
	for name := range buckets {
		names = append(names, name)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, name := range names {
		groups = append(groups, Group{Name: name, Keys: buckets[name]})
	}
	return groups
}
