// Package envcheck provides rule-based validation for environment variable maps.
//
// Rules can declare keys as required, forbidden, or non-empty. After running
// Check against an env map, any violations are returned as a slice that can be
// inspected with HasViolations and rendered with FormatViolations.
//
// Example:
//
//	rules := []envcheck.Rule{
//		{Key: "DATABASE_URL", Required: true, NonEmpty: true},
//		{Key: "DEBUG",        Forbidden: true},
//	}
//
//	violations := envcheck.Check(env, rules)
//	if envcheck.HasViolations(violations) {
//		for _, line := range envcheck.FormatViolations(violations) {
//			fmt.Println(line)
//		}
//	}
package envcheck
