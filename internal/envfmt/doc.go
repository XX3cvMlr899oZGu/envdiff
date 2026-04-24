// Package envfmt provides formatting utilities for env maps.
//
// It renders a map[string]string into canonical .env file format,
// with support for:
//
//   - Alphabetical key sorting
//   - Double or single quote wrapping of values
//   - Inline comments per key
//   - Custom key=value separator
//
// Example:
//
//	opts := envfmt.DefaultOptions()
//	opts.QuoteStyle = envfmt.QuoteDouble
//	opts.Comments = map[string]string{"PORT": "HTTP listen port"}
//	envfmt.Apply(os.Stdout, env, opts)
package envfmt
