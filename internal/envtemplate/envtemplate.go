// Package envtemplate renders Go text/template strings using env map values.
package envtemplate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		MissingKey: "error",
	}
}

// Options controls template rendering behaviour.
type Options struct {
	// MissingKey controls how missing keys are handled: "error", "zero", or "default".
	MissingKey string
}

// Result holds the rendered output for a single template string.
type Result struct {
	Input    string
	Rendered string
	Err      error
}

// Apply renders each value in tmplMap as a Go template, substituting variables
// from env. Keys in tmplMap that are not templates are passed through unchanged.
func Apply(tmplMap map[string]string, env map[string]string, opts Options) (map[string]string, error) {
	missingKey := opts.MissingKey
	if missingKey == "" {
		missingKey = "error"
	}

	out := make(map[string]string, len(tmplMap))
	for k, v := range tmplMap {
		if !strings.Contains(v, "{{}") && !strings.Contains(v, "{{") {
			out[k] = v
			continue
		}
		rendered, err := renderOne(v, env, missingKey)
		if err != nil {
			return nil, fmt.Errorf("envtemplate: key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}

// RenderString renders a single template string against env.
func RenderString(tmpl string, env map[string]string, opts Options) (string, error) {
	missingKey := opts.MissingKey
	if missingKey == "" {
		missingKey = "error"
	}
	return renderOne(tmpl, env, missingKey)
}

func renderOne(tmpl string, env map[string]string, missingKey string) (string, error) {
	t, err := template.New("").Option("missingkey=" + missingKey).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}
	var buf bytes.Buffer
	data := make(map[string]interface{}, len(env))
	for k, v := range env {
		data[k] = v
	}
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute error: %w", err)
	}
	return buf.String(), nil
}
