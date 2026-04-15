package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the export format type.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
)

// Write exports diff results to the given writer in the specified format.
func Write(w io.Writer, results []diff.Result, format Format) error {
	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	switch format {
	case FormatCSV:
		return writeCSV(w, sorted)
	case FormatMarkdown:
		return writeMarkdown(w, sorted)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func writeCSV(w io.Writer, results []diff.Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "status", "value_a", "value_b"}); err != nil {
		return err
	}
	for _, r := range results {
		if err := cw.Write([]string{r.Key, string(r.Status), r.ValueA, r.ValueB}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeMarkdown(w io.Writer, results []diff.Result) error {
	_, err := fmt.Fprintln(w, "| Key | Status | Value A | Value B |")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "|-----|--------|---------|---------|")
	if err != nil {
		return err
	}
	for _, r := range results {
		_, err = fmt.Fprintf(w, "| %s | %s | %s | %s |\n", r.Key, r.Status, r.ValueA, r.ValueB)
		if err != nil {
			return err
		}
	}
	return nil
}
