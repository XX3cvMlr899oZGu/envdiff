// Package envсummary provides a high-level summary of differences between
// two env maps, including counts by status and an overall health score.
package envsummary

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Summary holds aggregated statistics about a diff result.
type Summary struct {
	Total    int
	Equal    int
	Missing  int // present in A, absent in B
	Extra    int // present in B, absent in A
	Changed  int
	Score    float64 // 0.0 (all different) to 1.0 (identical)
	TopDiffs []string // up to N keys with differences, sorted
}

// DefaultOptions returns a Options with sensible defaults.
type Options struct {
	TopN int // max number of differing keys to include in TopDiffs
}

func DefaultOptions() Options {
	return Options{TopN: 5}
}

// Compute produces a Summary from two env maps using the diff package.
func Compute(a, b map[string]string, opts Options) Summary {
	results := diff.Compare(a, b)

	s := Summary{Total: len(results)}
	var diffKeys []string

	for _, r := range results {
		switch r.Status {
		case "equal":
			s.Equal++
		case "missing_in_second":
			s.Missing++
			diffKeys = append(diffKeys, r.Key)
		case "missing_in_first":
			s.Extra++
			diffKeys = append(diffKeys, r.Key)
		case "mismatch":
			s.Changed++
			diffKeys = append(diffKeys, r.Key)
		}
	}

	if s.Total > 0 {
		s.Score = float64(s.Equal) / float64(s.Total)
	} else {
		s.Score = 1.0
	}

	sort.Strings(diffKeys)
	n := opts.TopN
	if n > len(diffKeys) {
		n = len(diffKeys)
	}
	s.TopDiffs = diffKeys[:n]

	return s
}

// FormatText returns a human-readable summary string.
func FormatText(s Summary) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Total keys : %d\n", s.Total)
	fmt.Fprintf(&sb, "Equal      : %d\n", s.Equal)
	fmt.Fprintf(&sb, "Missing    : %d\n", s.Missing)
	fmt.Fprintf(&sb, "Extra      : %d\n", s.Extra)
	fmt.Fprintf(&sb, "Changed    : %d\n", s.Changed)
	fmt.Fprintf(&sb, "Score      : %.2f\n", s.Score)
	if len(s.TopDiffs) > 0 {
		fmt.Fprintf(&sb, "Top diffs  : %s\n", strings.Join(s.TopDiffs, ", "))
	}
	return sb.String()
}
