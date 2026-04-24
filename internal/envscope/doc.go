// Package envscope provides scoped views of environment variable maps.
//
// A Scope is defined by a name and a key prefix. Apply extracts all keys
// matching that prefix, optionally stripping the prefix from the resulting
// keys. ApplyAll handles multiple scopes at once, returning a map keyed by
// scope name.
//
// Example usage:
//
//	scope := envscope.Scope{Name: "db", Prefix: "DB_"}
//	result, err := envscope.Apply(env, scope, envscope.Options{StripPrefix: true})
//	// result: {"HOST": "localhost", "PORT": "5432"}
package envscope
