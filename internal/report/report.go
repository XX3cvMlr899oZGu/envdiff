package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Options configures report output.
type Options struct {
	Format    Format
	ShowEqual bool
}

// Write renders the diff results to the given writer.
func Write(w io.Writer, results []diff.Result, fileA, fileB string, opts Options) error {
	switch opts.Format {
	case FormatJSON:
		return writeJSON(w, results, fileA, fileB)
	default:
		return writeText(w, results, fileA, fileB, opts.ShowEqual)
	}
}

func writeText(w io.Writer, results []diff.Result, fileA, fileB string, showEqual bool) error {
	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(w, "Comparing: %s  →  %s\n", fileA, fileB)
	fmt.Fprintln(w, strings.Repeat("-", 48))

	for _, r := range sorted {
		switch r.Status {
		case diff.StatusMissingInA:
			fmt.Fprintf(w, "[+] %-30s  (only in %s)\n", r.Key, fileB)
		case diff.StatusMissingInB:
			fmt.Fprintf(w, "[-] %-30s  (only in %s)\n", r.Key, fileA)
		case diff.StatusMismatch:
			fmt.Fprintf(w, "[~] %-30s  %q → %q\n", r.Key, r.ValueA, r.ValueB)
		case diff.StatusEqual:
			if showEqual {
				fmt.Fprintf(w, "[=] %-30s\n", r.Key)
			}
		}
	}
	return nil
}
