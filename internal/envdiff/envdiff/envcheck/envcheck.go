// Package envcheck provides rule-based validation of env maps,
// checking for required keys, forbidden keys, and non-empty constraints.
package envcheck

import (
	"fmt"
	"strings"
)

// RuleKind defines the type of check to perform.
type RuleKind string

const (
	Required  RuleKind = "required"
	Forbidden RuleKind = "forbidden"
	NonEmpty  RuleKind = "non_empty"
)

// Rule describes a single validation rule.
type Rule struct {
	Key  string
	Kind RuleKind
}

// Violation describes a rule that was not satisfied.
type Violation struct {
	Key     string
	Kind    RuleKind
	Message string
}

// Check applies the given rules against env and returns any violations.
func Check(env map[string]string, rules []Rule) []Violation {
	var violations []Violation
	for _, r := range rules {
		switch r.Kind {
		case Required:
			if _, ok := env[r.Key]; !ok {
				violations = append(violations, Violation{
					Key:     r.Key,
					Kind:    r.Kind,
					Message: fmt.Sprintf("required key %q is missing", r.Key),
				})
			}
		case Forbidden:
			if _, ok := env[r.Key]; ok {
				violations = append(violations, Violation{
					Key:     r.Key,
					Kind:    r.Kind,
					Message: fmt.Sprintf("forbidden key %q is present", r.Key),
				})
			}
		case NonEmpty:
			if v, ok := env[r.Key]; !ok || strings.TrimSpace(v) == "" {
				violations = append(violations, Violation{
					Key:     r.Key,
					Kind:    r.Kind,
					Message: fmt.Sprintf("key %q must be present and non-empty", r.Key),
				})
			}
		}
	}
	return violations
}

// HasViolations returns true if any violations exist.
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
