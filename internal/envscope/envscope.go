// Package envscope provides scoped views of environment maps,
// allowing callers to extract a subset of keys matching a named scope.
package envscope

import (
	"fmt"
	"strings"
)

// Scope represents a named subset of environment keys.
type Scope struct {
	Name   string
	Prefix string
}

// Options controls how scopes are applied.
type Options struct {
	// StripPrefix removes the scope prefix from keys in the result.
	StripPrefix bool
	// CaseFold performs case-insensitive prefix matching.
	CaseFold bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		StripPrefix: false,
		CaseFold:    false,
	}
}

// Apply returns a new map containing only keys that belong to the given scope.
// If opts.StripPrefix is true, the scope prefix is removed from each key.
func Apply(env map[string]string, scope Scope, opts Options) (map[string]string, error) {
	if scope.Prefix == "" {
		return nil, fmt.Errorf("envscope: scope %q has empty prefix", scope.Name)
	}

	result := make(map[string]string)
	for k, v := range env {
		key := k
		prefix := scope.Prefix
		if opts.CaseFold {
			key = strings.ToLower(k)
			prefix = strings.ToLower(scope.Prefix)
		}
		if strings.HasPrefix(key, prefix) {
			outKey := k
			if opts.StripPrefix {
				outKey = k[len(scope.Prefix):]
				if outKey == "" {
					continue
				}
			}
			result[outKey] = v
		}
	}
	return result, nil
}

// ApplyAll applies multiple scopes and returns a map of scope name → scoped env.
func ApplyAll(env map[string]string, scopes []Scope, opts Options) (map[string]map[string]string, error) {
	out := make(map[string]map[string]string, len(scopes))
	for _, s := range scopes {
		scoped, err := Apply(env, s, opts)
		if err != nil {
			return nil, err
		}
		out[s.Name] = scoped
	}
	return out, nil
}
