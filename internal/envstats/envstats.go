// Package envstats provides statistical analysis of env maps,
// including value length distribution, key naming conventions, and
// basic summary metrics.
package envstats

import (
	"math"
	"sort"
	"strings"
)

// Stats holds summary statistics for an env map.
type Stats struct {
	TotalKeys      int
	EmptyValues    int
	AvgValueLength float64
	MaxValueLength int
	MinValueLength int
	UppercaseKeys  int
	LowercaseKeys  int
	MixedCaseKeys  int
	UniqueValues   int
}

// Compute calculates statistics from the given env map.
func Compute(env map[string]string) Stats {
	if len(env) == 0 {
		return Stats{}
	}

	total := len(env)
	empty := 0
	upper := 0
	lower := 0
	mixed := 0
	totalLen := 0
	maxLen := 0
	minLen := math.MaxInt64
	seen := make(map[string]struct{})

	for k, v := range env {
		vLen := len(v)
		if v == "" {
			empty++
		}
		totalLen += vLen
		if vLen > maxLen {
			maxLen = vLen
		}
		if vLen < minLen {
			minLen = vLen
		}
		seen[v] = struct{}{}

		kUpper := strings.ToUpper(k)
		kLower := strings.ToLower(k)
		switch {
		case k == kUpper:
			upper++
		case k == kLower:
			lower++
		default:
			mixed++
		}
	}

	if minLen == math.MaxInt64 {
		minLen = 0
	}

	return Stats{
		TotalKeys:      total,
		EmptyValues:    empty,
		AvgValueLength: float64(totalLen) / float64(total),
		MaxValueLength: maxLen,
		MinValueLength: minLen,
		UppercaseKeys:  upper,
		LowercaseKeys:  lower,
		MixedCaseKeys:  mixed,
		UniqueValues:   len(seen),
	}
}

// TopLongestValues returns the top n keys by value length, sorted descending.
func TopLongestValues(env map[string]string, n int) []string {
	type kv struct {
		key string
		len int
	}
	pairs := make([]kv, 0, len(env))
	for k, v := range env {
		pairs = append(pairs, kv{k, len(v)})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].len != pairs[j].len {
			return pairs[i].len > pairs[j].len
		}
		return pairs[i].key < pairs[j].key
	})
	result := make([]string, 0, n)
	for i := 0; i < n && i < len(pairs); i++ {
		result = append(result, pairs[i].key)
	}
	return result
}
