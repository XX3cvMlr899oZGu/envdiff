// Package envcheck provides rule-based validation of environment variable maps.
// It supports required keys, forbidden keys, and non-empty value constraints.
package envcheck

import (
	"fmt"
	"strings"
)

// RuleType defines the kind of constraint to enforce.
type RuleType string

const (
	RuleRequired RuleType = "required"
	RuleForbidden RuleType = "forbidden"
	RuleNonEmpty  RuleType = "non_empty"
)

// Rule describes a single constraint applied to an env map.
type Rule struct {
	Key  string
	Type RuleType
}

// Violation records a rule that was not satisfied.
type Violation struct {
	Key     string
	Rule    RuleType
	Message string
}

// Check evaluates a set of rules against the provided env map and returns
// any violations found.
func Check(env map[string]string, rules []Rule) []Violation {
	var violations []Violation
	for _, r := range rules {
		switch r.Type {
		case RuleRequired:
			if _, ok := env[r.Key]; !ok {
				violations = append(violations, Violation{
					Key:     r.Key,
					Rule:    r.Type,
					Message: fmt.Sprintf("required key %q is missing", r.Key),
				})
			}
		case RuleForbidden:
			if _, ok := env[r.Key]; ok {
				violations = append(violations, Violation{
					Key:     r.Key,
					Rule:    r.Type,
					Message: fmt.Sprintf("forbidden key %q is present", r.Key),
				})
			}
		case RuleNonEmpty:
			if v, ok := env[r.Key]; ok && strings.TrimSpace(v) == "" {
				violations = append(violations, Violation{
					Key:     r.Key,
					Rule:    r.Type,
					Message: fmt.Sprintf("key %q must not be empty", r.Key),
				})
			}
		}
	}
	return violations
}

// HasViolations returns true if the slice contains at least one violation.
func HasViolations(vs []Violation) bool {
	return len(vs) > 0
}

// FormatViolations returns a human-readable summary of all violations.
func FormatViolations(vs []Violation) string {
	if len(vs) == 0 {
		return "no violations"
	}
	lines := make([]string, len(vs))
	for i, v := range vs {
		lines[i] = fmt.Sprintf("[%s] %s", v.Rule, v.Message)
	}
	return strings.Join(lines, "\n")
}
