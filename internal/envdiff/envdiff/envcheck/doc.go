// Package envcheck provides a lightweight rule engine for validating
// environment variable maps.
//
// Supported rule kinds:
//
//   - required  — the key must be present in the env map
//   - forbidden — the key must NOT be present in the env map
//   - non_empty — if the key is present, its value must not be blank
//
// Example usage:
//
//	rules := []envcheck.Rule{
//		{Key: "DATABASE_URL", Kind: envcheck.RuleRequired},
//		{Key: "DEBUG",        Kind: envcheck.RuleForbidden},
//		{Key: "API_KEY",      Kind: envcheck.RuleNonEmpty},
//	}
//
//	violations := envcheck.Check(env, rules)
//	if envcheck.HasViolations(violations) {
//		fmt.Println(envcheck.FormatViolations(violations))
//	}
package envcheck
