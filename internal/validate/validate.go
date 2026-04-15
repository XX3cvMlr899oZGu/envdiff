package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for environment variable keys or values.
type Rule struct {
	Key     string
	Pattern string
	Required bool
}

// Violation represents a failed validation rule.
type Violation struct {
	Key     string
	Message string
}

// ApplyRules validates an env map against a set of rules and returns any violations.
func ApplyRules(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, val)
			if err != nil {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !matched {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
			}
		}
	}

	return violations
}

// FormatViolations returns a human-readable summary of violations.
func FormatViolations(violations []Violation) string {
	if len(violations) == 0 {
		return "no validation violations found"
	}
	var sb strings.Builder
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  [VIOLATION] %s: %s\n", v.Key, v.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}
