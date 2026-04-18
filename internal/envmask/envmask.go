// Package envmask provides masking of env values based on key patterns,
// replacing a portion of the value with asterisks to partially obscure it.
package envmask

import (
	"regexp"
	"strings"
)

// Options controls masking behaviour.
type Options struct {
	// Patterns is a list of regexp patterns; matching keys will be masked.
	Patterns []string
	// VisibleChars is how many leading chars to keep visible (default 2).
	VisibleChars int
	// Placeholder replaces the hidden portion (default "****").
	Placeholder string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Patterns:     []string{`(?i)(secret|token|password|key|pwd|auth)`},
		VisibleChars: 2,
		Placeholder:  "****",
	}
}

// Apply returns a new map with sensitive values partially masked.
func Apply(env map[string]string, opts Options) (map[string]string, error) {
	if opts.Placeholder == "" {
		opts.Placeholder = "****"
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
			out[k] = maskValue(v, opts.VisibleChars, opts.Placeholder)
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func isSensitive(key string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

func maskValue(v string, visible int, placeholder string) string {
	if len(v) == 0 {
		return placeholder
	}
	if visible <= 0 || visible >= len(v) {
		return placeholder
	}
	return v[:visible] + strings.Repeat("*", len(placeholder))
}
