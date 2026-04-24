// Package envcheck provides rule-based validation checks against an env map.
// Rules can require keys to be present, forbid keys, or enforce non-empty values.
package envcheck

import (
	"fmt"
	"strings"
)

// RuleKind defines the type of check to perform.
type RuleKind string

const (
	RuleRequired RuleKind = "required"
	RuleForbidden RuleKind = "forbidden"
	RuleNonEmpty  RuleKind = "non_empty"
)

// Rule describes a single check against an env map.
type Rule struct {
	Key  string
	Kind RuleKind
}

// Violation represents a failed rule check.
type Violation struct {
	Key     string
	Kind    RuleKind
	Message string
}

// Check applies the given rules to env and returns any violations.
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
			val, ok := env[r.Key]
			if ok && strings.TrimSpace(val) == "" {
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

// FormatViolations returns a human-readable summary of violations.
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
