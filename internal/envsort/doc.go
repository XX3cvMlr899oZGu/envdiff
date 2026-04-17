// Package envsort provides deterministic ordering of environment variable maps.
//
// Env maps in Go are unordered by nature. envsort converts them to sorted
// slices of Entry values, which is useful for producing stable output in
// reports, exports, and diffs.
//
// Sorting options include:
//   - Ascending / Descending alphabetical order
//   - Prefix priority: keys matching a given prefix are sorted before others
//
// Example:
//
//	env := map[string]string{"Z": "1", "A": "2"}
//	entries := envsort.Apply(env, envsort.DefaultOptions())
//	// entries[0].Key == "A", entries[1].Key == "Z"
package envsort
