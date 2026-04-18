// Package envdefault fills in missing keys from a defaults map.
package envdefault

// Options controls how defaults are applied.
type Options struct {
	// Overwrite replaces existing values with defaults when true.
	Overwrite bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Overwrite: false}
}

// Apply merges defaults into env. Returns a new map; env is not mutated.
// Keys present in defaults but absent in env are added.
// If opts.Overwrite is true, all keys from defaults overwrite env values.
func Apply(env map[string]string, defaults map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for k, v := range defaults {
		existing, exists := out[k]
		_ = existing
		if !exists || opts.Overwrite {
			out[k] = v
		}
	}
	return out
}

// MissingKeys returns keys that are in defaults but not in env.
func MissingKeys(env map[string]string, defaults map[string]string) []string {
	var missing []string
	for k := range defaults {
		if _, ok := env[k]; !ok {
			missing = append(missing, k)
		}
	}
	return missing
}
