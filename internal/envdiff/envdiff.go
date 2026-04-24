// Package envdiff provides summary statistics over a set of diff results.
package envdiff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Summary holds aggregate counts over a slice of diff results.
type Summary struct {
	Total    int
	Equal    int
	Missing  int
	Extra    int
	Mismatch int
}

// Summarize computes a Summary from a slice of diff.Result values.
func Summarize(results []diff.Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case diff.Equal:
			s.Equal++
		case diff.MissingInSecond:
			s.Missing++
		case diff.MissingInFirst:
			s.Extra++
		case diff.ValueMismatch:
			s.Mismatch++
		}
	}
	return s
}

// HasDifferences returns true when any result is not Equal.
func HasDifferences(results []diff.Result) bool {
	for _, r := range results {
		if r.Status != diff.Equal {
			return true
		}
	}
	return false
}

// FormatSummary returns a human-readable one-line summary string.
func FormatSummary(s Summary) string {
	parts := []string{
		fmt.Sprintf("total=%d", s.Total),
		fmt.Sprintf("equal=%d", s.Equal),
		fmt.Sprintf("missing=%d", s.Missing),
		fmt.Sprintf("extra=%d", s.Extra),
		fmt.Sprintf("mismatch=%d", s.Mismatch),
	}
	return strings.Join(parts, " ")
}

// KeysByStatus returns the sorted list of keys that have the given status.
func KeysByStatus(results []diff.Result, status diff.Status) []string {
	var keys []string
	for _, r := range results {
		if r.Status == status {
			keys = append(keys, r.Key)
		}
	}
	sort.Strings(keys)
	return keys
}
