// Package cast provides type coercion for environment variable values loaded
// from .env files. It converts raw string values into typed Go values (bool,
// int, float64, or string) based on caller-supplied type hints.
//
// # Usage
//
//	opts := cast.DefaultOptions()
//	opts.TypeHints = map[string]string{
//		"PORT":  "int",
//		"DEBUG": "bool",
//	}
//
//	results, err := cast.Apply(env, opts)
//
// Keys without a type hint are returned as strings when FallbackToString is
// true (the default). Set FallbackToString to false to skip unrecognised keys
// entirely.
//
// Supported type names (case-insensitive):
//
//   - "string"  — identity conversion
//   - "int"     — parsed with strconv.Atoi
//   - "float"   — parsed with strconv.ParseFloat (64-bit)
//   - "bool"    — parsed with strconv.ParseBool
package cast
