package report

import (
	"encoding/json"
	"io"

	"github.com/user/envdiff/internal/diff"
)

type jsonReport struct {
	FileA   string       `json:"file_a"`
	FileB   string       `json:"file_b"`
	Results []jsonResult `json:"results"`
}

type jsonResult struct {
	Key     string `json:"key"`
	Status  string `json:"status"`
	ValueA  string `json:"value_a,omitempty"`
	ValueB  string `json:"value_b,omitempty"`
}

func writeJSON(w io.Writer, results []diff.Result, fileA, fileB string) error {
	jr := jsonReport{
		FileA:   fileA,
		FileB:   fileB,
		Results: make([]jsonResult, 0, len(results)),
	}

	for _, r := range results {
		jr.Results = append(jr.Results, jsonResult{
			Key:    r.Key,
			Status: statusString(r.Status),
			ValueA: r.ValueA,
			ValueB: r.ValueB,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jr)
}

func statusString(s diff.Status) string {
	switch s {
	case diff.StatusMissingInA:
		return "missing_in_a"
	case diff.StatusMissingInB:
		return "missing_in_b"
	case diff.StatusMismatch:
		return "mismatch"
	case diff.StatusEqual:
		return "equal"
	default:
		return "unknown"
	}
}
