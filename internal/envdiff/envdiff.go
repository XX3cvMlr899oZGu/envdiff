// Package envdiff provides utilities for summarizing diff statistics.
package envdiff

import "github.com/yourorg/envdiff/internal/diff"

// Stats holds a summary of comparison results.
type Stats struct {
	Total     int
	Missing   int
	Extra     int
	Mismatch  int
	Equal     int
}

// Summarize computes statistics from a slice of diff.Result.
func Summarize(results []diff.Result) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case diff.Missing:
			s.Missing++
		case diff.Extra:
			s.Extra++
		case diff.Mismatch:
			s.Mismatch++
		case diff.Equal:
			s.Equal++
		}
	}
	return s
}

// HasDifferences returns true if any result is not Equal.
func HasDifferences(results []diff.Result) bool {
	for _, r := range results {
		if r.Status != diff.Equal {
			return true
		}
	}
	return false
}
