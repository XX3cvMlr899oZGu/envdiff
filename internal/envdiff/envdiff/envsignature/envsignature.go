// Package envsignature generates a stable, human-readable signature string
// for an env map, useful for quick identity checks across environments.
package envsignature

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Options controls signature generation.
type Options struct {
	// IncludeValues includes key values in the signature when true.
	// When false, only keys are hashed (useful for structure-only comparison).
	IncludeValues bool

	// Prefix restricts the signature to keys with the given prefix.
	Prefix string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		IncludeValues: true,
	}
}

// Result holds the computed signature and metadata.
type Result struct {
	// Signature is a short hex string derived from the env map.
	Signature string
	// KeyCount is the number of keys that contributed to the signature.
	KeyCount int
}

// Compute generates a deterministic signature for the given env map.
func Compute(env map[string]string, opts Options) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		if opts.IncludeValues {
			fmt.Fprintf(h, "%s=%s\n", k, env[k])
		} else {
			fmt.Fprintf(h, "%s\n", k)
		}
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return Result{
		Signature: sum[:16],
		KeyCount:  len(keys),
	}
}

// Equal returns true if two env maps produce the same signature under opts.
func Equal(a, b map[string]string, opts Options) bool {
	return Compute(a, opts).Signature == Compute(b, opts).Signature
}
