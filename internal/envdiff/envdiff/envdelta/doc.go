// Package envdelta computes a structured, key-level delta between two env maps.
//
// Each key is classified as one of:
//   - added:     present in next but not in base
//   - removed:   present in base but not in next
//   - changed:   present in both but with different values
//   - unchanged: present in both with identical values
//
// Usage:
//
//	delta := envdelta.Compute(baseEnv, nextEnv, envdelta.DefaultOptions())
//	if delta.HasChanges() {
//	    for _, e := range delta.ByStatus(envdelta.StatusChanged) {
//	        fmt.Printf("~ %s: %q -> %q\n", e.Key, e.OldVal, e.NewVal)
//	    }
//	}
package envdelta
