// Package envfreeze provides functionality to "freeze" an env map into a
// locked snapshot and detect any deviation from it later.
package envfreeze

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// FreezeRecord holds a frozen snapshot of an env map with metadata.
type FreezeRecord struct {
	FrozenAt time.Time         `json:"frozen_at"`
	Env      map[string]string `json:"env"`
}

// DefaultOptions returns a zero-value Options suitable for most callers.
func DefaultOptions() Options {
	return Options{}
}

// Options controls freeze/thaw behaviour.
type Options struct {
	// IgnoreKeys are keys excluded from deviation checks.
	IgnoreKeys []string
}

// Freeze writes a FreezeRecord for env to path.
func Freeze(env map[string]string, path string) error {
	record := FreezeRecord{
		FrozenAt: time.Now().UTC(),
		Env:      env,
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("envfreeze: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("envfreeze: write %s: %w", path, err)
	}
	return nil
}

// Load reads a FreezeRecord from path.
func Load(path string) (FreezeRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FreezeRecord{}, fmt.Errorf("envfreeze: read %s: %w", path, err)
	}
	var record FreezeRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return FreezeRecord{}, fmt.Errorf("envfreeze: unmarshal: %w", err)
	}
	return record, nil
}

// Deviation describes a single key that differs from the frozen state.
type Deviation struct {
	Key      string
	Kind     string // "added", "removed", "changed"
	Frozen   string
	Current  string
}

// Diff compares current env against a FreezeRecord and returns deviations.
func Diff(record FreezeRecord, current map[string]string, opts Options) []Deviation {
	ignore := make(map[string]bool, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignore[k] = true
	}

	var deviations []Deviation

	for k, fv := range record.Env {
		if ignore[k] {
			continue
		}
		cv, ok := current[k]
		if !ok {
			deviations = append(deviations, Deviation{Key: k, Kind: "removed", Frozen: fv})
		} else if cv != fv {
			deviations = append(deviations, Deviation{Key: k, Kind: "changed", Frozen: fv, Current: cv})
		}
	}

	for k, cv := range current {
		if ignore[k] {
			continue
		}
		if _, ok := record.Env[k]; !ok {
			deviations = append(deviations, Deviation{Key: k, Kind: "added", Current: cv})
		}
	}

	sort.Slice(deviations, func(i, j int) bool {
		return deviations[i].Key < deviations[j].Key
	})
	return deviations
}

// HasDeviations returns true when Diff produces at least one entry.
func HasDeviations(deviations []Deviation) bool { return len(deviations) > 0 }

// FormatText formats deviations as a human-readable string.
func FormatText(deviations []Deviation) string {
	if len(deviations) == 0 {
		return "no deviations from frozen state\n"
	}
	out := ""
	for _, d := range deviations {
		switch d.Kind {
		case "added":
			out += fmt.Sprintf("+ %s = %q (not in freeze)\n", d.Key, d.Current)
		case "removed":
			out += fmt.Sprintf("- %s (was %q, now missing)\n", d.Key, d.Frozen)
		case "changed":
			out += fmt.Sprintf("~ %s: frozen=%q current=%q\n", d.Key, d.Frozen, d.Current)
		}
	}
	return out
}
