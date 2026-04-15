package redact

import (
	"regexp"
	"strings"
)

// DefaultPatterns contains common patterns for sensitive keys.
var DefaultPatterns = []string{
	"(?i)password",
	"(?i)secret",
	"(?i)token",
	"(?i)api_key",
	"(?i)private_key",
	"(?i)auth",
}

// Options controls redaction behaviour.
type Options struct {
	// Patterns is a list of regex patterns matched against key names.
	Patterns []string
	// Placeholder is the string used to replace sensitive values.
	Placeholder string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Patterns:    DefaultPatterns,
		Placeholder: "***REDACTED***",
	}
}

// Apply returns a copy of env with sensitive values replaced by the placeholder.
func Apply(env map[string]string, opts Options) (map[string]string, error) {
	if opts.Placeholder == "" {
		opts.Placeholder = "***REDACTED***"
	}

	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, compiled) {
			out[k] = opts.Placeholder
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// IsSensitiveKey reports whether a key matches any of the default patterns.
func IsSensitiveKey(key string) bool {
	for _, p := range DefaultPatterns {
		re := regexp.MustCompile(p)
		if re.MatchString(strings.ToLower(key)) {
			return true
		}
	}
	return false
}

func isSensitive(key string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
