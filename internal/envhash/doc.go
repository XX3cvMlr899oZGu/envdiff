// Package envhash provides deterministic SHA-256 hashing of env maps.
//
// It is useful for detecting whether an environment has changed between
// runs, CI stages, or deployments without comparing individual keys.
//
// Basic usage:
//
//	hash, err := envhash.Compute(env, envhash.DefaultOptions())
//
// Keys are sorted before hashing to guarantee that insertion order does
// not affect the result. Sensitive keys can be excluded via Options.ExcludeKeys,
// and hashing can be scoped to a subset of keys via Options.IncludeKeys.
//
// Use Equal for a convenient two-map comparison:
//
//	same, err := envhash.Equal(envA, envB, envhash.DefaultOptions())
package envhash
