// Package envlabel attaches arbitrary string labels to environment maps,
// enabling downstream tooling to annotate, filter, or route env sets by tag.
package envlabel

import (
	"errors"
	"fmt"
	"sort"
)

// DefaultOptions returns a zero-value Options ready for use.
func DefaultOptions() Options {
	return Options{}
}

// Options controls labelling behaviour.
type Options struct {
	// Labels is the set of key/value labels to attach.
	Labels map[string]string
	// OverwriteExisting allows existing label keys to be replaced.
	OverwriteExisting bool
}

// LabeledEnv pairs an environment map with its labels.
type LabeledEnv struct {
	Env    map[string]string
	Labels map[string]string
}

// Apply attaches labels from opts to env, returning a LabeledEnv.
// If OverwriteExisting is false and a label key already exists, an error is returned.
func Apply(env map[string]string, existing map[string]string, opts Options) (LabeledEnv, error) {
	if env == nil {
		return LabeledEnv{}, errors.New("envlabel: env must not be nil")
	}

	out := make(map[string]string, len(existing))
	for k, v := range existing {
		out[k] = v
	}

	for k, v := range opts.Labels {
		if k == "" {
			return LabeledEnv{}, errors.New("envlabel: label key must not be empty")
		}
		if _, exists := out[k]; exists && !opts.OverwriteExisting {
			return LabeledEnv{}, fmt.Errorf("envlabel: label key %q already exists", k)
		}
		out[k] = v
	}

	envCopy := make(map[string]string, len(env))
	for k, v := range env {
		envCopy[k] = v
	}

	return LabeledEnv{Env: envCopy, Labels: out}, nil
}

// FormatText returns a human-readable summary of the labels.
func FormatText(le LabeledEnv) string {
	if len(le.Labels) == 0 {
		return "(no labels)"
	}
	keys := make([]string, 0, len(le.Labels))
	for k := range le.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := ""
	for _, k := range keys {
		out += fmt.Sprintf("%s=%s\n", k, le.Labels[k])
	}
	return out
}
