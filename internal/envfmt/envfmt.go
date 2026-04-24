// Package envfmt formats env maps into canonical .env file output.
// It supports key sorting, section separators, and comment injection.
package envfmt

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		SortKeys:   true,
		EqualSign:  "=",
		QuoteStyle: QuoteNone,
	}
}

// QuoteStyle controls how values are quoted in output.
type QuoteStyle int

const (
	QuoteNone   QuoteStyle = iota // no quoting
	QuoteDouble                   // wrap all values in double quotes
	QuoteSingle                   // wrap all values in single quotes
)

// Options configures the formatter.
type Options struct {
	SortKeys   bool
	EqualSign  string
	QuoteStyle QuoteStyle
	Comments   map[string]string // key -> inline comment
}

// Apply writes the env map to w using the given options.
func Apply(w io.Writer, env map[string]string, opts Options) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	eq := opts.EqualSign
	if eq == "" {
		eq = "="
	}

	for _, k := range keys {
		v := quoteValue(env[k], opts.QuoteStyle)
		line := fmt.Sprintf("%s%s%s", k, eq, v)
		if c, ok := opts.Comments[k]; ok && c != "" {
			line += " # " + strings.TrimSpace(c)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func quoteValue(v string, style QuoteStyle) string {
	switch style {
	case QuoteDouble:
		return `"` + v + `"`
	case QuoteSingle:
		return `'` + v + `'`
	default:
		return v
	}
}
