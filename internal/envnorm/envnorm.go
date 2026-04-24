// Package envnorm provides key normalization for env maps.
// It supports converting keys to uppercase, lowercase, or snake_case,
// and optionally replacing hyphens or dots with underscores.
package envnorm

import (
	"fmt"
	"strings"
)

// Style defines the normalization style applied to keys.
type Style string

const (
	StyleUpper  Style = "upper"
	StyleLower  Style = "lower"
	StyleSnake  Style = "snake"
	StyleNone   Style = "none"
)

// Options controls how keys are normalized.
type Options struct {
	Style           Style
	ReplaceHyphens  bool // replace '-' with '_'
	ReplaceDots     bool // replace '.' with '_'
}

// DefaultOptions returns sensible defaults: uppercase, replace hyphens.
func DefaultOptions() Options {
	return Options{
		Style:          StyleUpper,
		ReplaceHyphens: true,
		ReplaceDots:    false,
	}
}

// Apply normalizes all keys in env according to opts.
// Returns an error if two distinct original keys normalize to the same result.
func Apply(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	seen := make(map[string]string, len(env)) // normalized -> original

	for k, v := range env {
		nk := normalizeKey(k, opts)
		if orig, exists := seen[nk]; exists && orig != k {
			return nil, fmt.Errorf("envnorm: key collision after normalization: %q and %q both become %q", orig, k, nk)
		}
		seen[nk] = k
		out[nk] = v
	}
	return out, nil
}

func normalizeKey(k string, opts Options) string {
	if opts.ReplaceHyphens {
		k = strings.ReplaceAll(k, "-", "_")
	}
	if opts.ReplaceDots {
		k = strings.ReplaceAll(k, ".", "_")
	}
	switch opts.Style {
	case StyleUpper:
		return strings.ToUpper(k)
	case StyleLower:
		return strings.ToLower(k)
	case StyleSnake:
		// snake: uppercase + underscore-separated (already handled above)
		return strings.ToUpper(k)
	default:
		return k
	}
}
