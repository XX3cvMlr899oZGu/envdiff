// Package envhash computes deterministic hashes of env maps for change detection.
package envhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Options controls hash computation behaviour.
type Options struct {
	// IncludeKeys limits hashing to these keys. If empty, all keys are used.
	IncludeKeys []string
	// ExcludeKeys are keys to skip when computing the hash.
	ExcludeKeys []string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{}
}

// Compute returns a hex-encoded SHA-256 hash of the env map.
// Keys are sorted before hashing to ensure determinism.
func Compute(env map[string]string, opts Options) (string, error) {
	exclude := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		exclude[k] = struct{}{}
	}

	include := make(map[string]struct{}, len(opts.IncludeKeys))
	for _, k := range opts.IncludeKeys {
		include[k] = struct{}{}
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		if _, skip := exclude[k]; skip {
			continue
		}
		if len(include) > 0 {
			if _, ok := include[k]; !ok {
				continue
			}
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		line := fmt.Sprintf("%s=%s\n", k, env[k])
		if _, err := h.Write([]byte(line)); err != nil {
			return "", fmt.Errorf("envhash: write error: %w", err)
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Equal returns true when two env maps produce the same hash.
func Equal(a, b map[string]string, opts Options) (bool, error) {
	ha, err := Compute(a, opts)
	if err != nil {
		return false, err
	}
	hb, err := Compute(b, opts)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(ha, hb), nil
}
