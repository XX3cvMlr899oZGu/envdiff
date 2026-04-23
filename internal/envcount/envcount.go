// Package envcount provides utilities for counting and summarizing
// key statistics about one or more parsed environment maps.
package envcount

import (
	"sort"
	"strings"
)

// Stats holds aggregate counts for an environment map.
type Stats struct {
	Total     int
	Empty     int
	NonEmpty  int
	Prefixes  map[string]int // count of keys per top-level prefix (e.g. "DB" for "DB_HOST")
}

// Options controls how counting is performed.
type Options struct {
	// Separator is the delimiter used to extract the prefix from a key.
	// Defaults to "_".
	Separator string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
	}
}

// Apply computes Stats for the given environment map using the provided options.
func Apply(env map[string]string, opts Options) Stats {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	s := Stats{
		Prefixes: make(map[string]int),
	}

	for k, v := range env {
		s.Total++
		if strings.TrimSpace(v) == "" {
			s.Empty++
		} else {
			s.NonEmpty++
		}

		parts := strings.SplitN(k, opts.Separator, 2)
		if len(parts) == 2 && parts[0] != "" {
			s.Prefixes[parts[0]]++
		}
	}

	return s
}

// TopPrefixes returns the top-n prefixes by key count, sorted descending.
// If n <= 0, all prefixes are returned.
func TopPrefixes(s Stats, n int) []string {
	type kv struct {
		key   string
		count int
	}

	pairs := make([]kv, 0, len(s.Prefixes))
	for k, c := range s.Prefixes {
		pairs = append(pairs, kv{k, c})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})

	result := make([]string, 0, len(pairs))
	for idx, p := range pairs {
		if n > 0 && idx >= n {
			break
		}
		result = append(result, p.key)
	}
	return result
}
