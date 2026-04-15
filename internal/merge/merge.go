package merge

import (
	"fmt"
	"sort"
)

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy int

const (
	// StrategyFirst keeps the value from the first (base) map on conflict.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last map on conflict.
	StrategyLast
	// StrategyError returns an error on any conflicting key.
	StrategyError
)

// Result holds the merged environment and any keys that had conflicts.
type Result struct {
	Env       map[string]string
	Conflicts []string
}

// Merge combines multiple env maps into one according to the given strategy.
// Maps are applied left-to-right; earlier maps are considered the "base".
func Merge(strategy Strategy, maps ...map[string]string) (*Result, error) {
	merged := make(map[string]string)
	conflictSet := make(map[string]bool)

	for _, m := range maps {
		for k, v := range m {
			existing, exists := merged[k]
			if !exists {
				merged[k] = v
				continue
			}
			if existing == v {
				continue
			}
			// Conflict detected
			switch strategy {
			case StrategyFirst:
				// keep existing, record conflict
				conflictSet[k] = true
			case StrategyLast:
				merged[k] = v
				conflictSet[k] = true
			case StrategyError:
				return nil, fmt.Errorf("merge conflict on key %q: %q vs %q", k, existing, v)
			}
		}
	}

	conflicts := make([]string, 0, len(conflictSet))
	for k := range conflictSet {
		conflicts = append(conflicts, k)
	}
	sort.Strings(conflicts)

	return &Result{Env: merged, Conflicts: conflicts}, nil
}
