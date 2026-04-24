// Package envcheck provides rule-based validation of environment variable maps,
// checking for required keys, forbidden keys, and non-empty constraints.
package envcheck

import (
	"fmt"
	"strings"
)

// RuleKind describes the type of constraint to enforce.
type RuleKind string

const (
	RuleRequired  RuleKind = "required"
	RuleForbidden RuleKind = "forbidden"
	RuleNonEmpty  RuleKind = "non_empty"
)

// Rule defines a single constraint applied to an env map.
type Rule struct {
	Key  string
	Kind RuleKind
}

// Violation describes a single failed rule check.
type Violation struct {
	Key     string
	Kind    RuleKind
	Message string
}

// Check applies all rules to env and returns any violations found.
func Check(env map[string]string, rules []Rule) []Violation {
	var violations []Violation
	for _, r := range rules {
		switch r.Kind {
		case RuleRequired:
			if _, ok := env[r.Key]; !ok {
				violations = append(violations, Violation{
					Key:     r.Key,
					Kind:    r.Kind,
					Message: fmt.Sprintf("required key %q is missing", r.Key),
				})
			}
		case RuleForbidden:
			if _, ok := env[r.Key]; ok {
				violations = append(violations, Violation{
					Key:     r.Key,
					Kind:    r.Kind,
					Message: fmt.Sprintf("forbidden key %q is present", r.Key),
				})
			}
		case RuleNonEmpty:
			if v, ok := env[r.Key]; ok && strings.TrimSpace(v) == "" {
				violations = append(violations, Violation{
					Key:     r.Key,
					Kind:    r.Kind,
					Message: fmt.Sprintf("key %q must not be empty", r.Key),
				})
			}
		}
	}
	return violations
}

// HasViolations returns true if any violations are present.
func HasViolations(violations []Violation) bool {
	return len(violations) > 0
}

// FormatViolations returns a human-readable summary of all violations.
func FormatViolations(violations []Violation) string {
	if len(violations) == 0 {
		return "no violations found"
	}
	var sb strings.Builder
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", v.Kind, v.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}
