// Package envcheck provides utilities for checking env maps against
// a set of required and forbidden key rules.
package envcheck

import (
	"fmt"
	"sort"
)

// Rule describes a single check to apply against an env map.
type Rule struct {
	// Key is the environment variable name to check.
	Key string
	// Required indicates the key must be present (and non-empty if NonEmpty is true).
	Required bool
	// Forbidden indicates the key must NOT be present.
	Forbidden bool
	// NonEmpty requires the value to be non-empty when the key is present.
	NonEmpty bool
}

// Violation describes a single failed check.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Check applies rules against the provided env map and returns any violations.
func Check(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, r := range rules {
		val, exists := env[r.Key]

		if r.Forbidden && exists {
			violations = append(violations, Violation{
				Key:     r.Key,
				Message: "key is forbidden but present",
			})
			continue
		}

		if r.Required && !exists {
			violations = append(violations, Violation{
				Key:     r.Key,
				Message: "required key is missing",
			})
			continue
		}

		if r.NonEmpty && exists && val == "" {
			violations = append(violations, Violation{
				Key:     r.Key,
				Message: "key is present but value is empty",
			})
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Key < violations[j].Key
	})

	return violations
}

// HasViolations returns true if Check produces any violations.
func HasViolations(env map[string]string, rules []Rule) bool {
	return len(Check(env, rules)) > 0
}

// FormatViolations returns a human-readable summary of violations.
func FormatViolations(violations []Violation) string {
	if len(violations) == 0 {
		return "no violations found"
	}
	result := fmt.Sprintf("%d violation(s):\n", len(violations))
	for _, v := range violations {
		result += fmt.Sprintf("  [%s] %s\n", v.Key, v.Message)
	}
	return result
}
