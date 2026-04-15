// Package lint provides heuristic checks on parsed .env maps,
// warning about common issues such as empty values, keys with
// whitespace, or suspiciously long values.
package lint

import (
	"fmt"
	"strings"
)

// Severity indicates how serious a lint warning is.
type Severity string

const (
	Warn  Severity = "WARN"
	Error Severity = "ERROR"
)

// Issue represents a single lint finding.
type Issue struct {
	Key      string
	Message  string
	Severity Severity
}

// Options controls which checks are enabled.
type Options struct {
	MaxValueLen    int  // 0 means no limit
	WarnEmptyValue bool
	WarnKeySpaces  bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxValueLen:    500,
		WarnEmptyValue: true,
		WarnKeySpaces:  true,
	}
}

// Run applies all enabled checks to env and returns any issues found.
func Run(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for k, v := range env {
		if opts.WarnKeySpaces && strings.ContainsAny(k, " \t") {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "key contains whitespace",
				Severity: Error,
			})
		}

		if opts.WarnEmptyValue && strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "value is empty",
				Severity: Warn,
			})
		}

		if opts.MaxValueLen > 0 && len(v) > opts.MaxValueLen {
			issues = append(issues, Issue{
				Key:      k,
				Message:  fmt.Sprintf("value length %d exceeds maximum %d", len(v), opts.MaxValueLen),
				Severity: Warn,
			})
		}
	}

	return issues
}

// FormatIssues returns a human-readable summary of the issues slice.
func FormatIssues(issues []Issue) string {
	if len(issues) == 0 {
		return "no lint issues found"
	}
	var sb strings.Builder
	for _, iss := range issues {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", iss.Severity, iss.Key, iss.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
