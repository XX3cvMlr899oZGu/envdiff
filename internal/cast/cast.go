// Package cast provides type coercion utilities for environment variable values.
// It converts string values from .env files into typed Go values.
package cast

import (
	"fmt"
	"strconv"
	"strings"
)

// Result holds a typed value parsed from an env string.
type Result struct {
	Key      string
	Raw      string
	TypeName string
	Value    interface{}
}

// Options controls how casting is performed.
type Options struct {
	// TypeHints maps key names to desired types: "bool", "int", "float", "string"
	TypeHints map[string]string
	// FallbackToString returns a string Result when no hint matches instead of an error
	FallbackToString bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		TypeHints:        map[string]string{},
		FallbackToString: true,
	}
}

// Apply casts each value in env according to opts.TypeHints.
// Keys without a hint are treated as strings when FallbackToString is true.
func Apply(env map[string]string, opts Options) ([]Result, error) {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		hint, ok := opts.TypeHints[k]
		if !ok {
			if opts.FallbackToString {
				results = append(results, Result{Key: k, Raw: v, TypeName: "string", Value: v})
			}
			continue
		}
		r, err := castValue(k, v, hint)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func castValue(key, raw, typeName string) (Result, error) {
	switch strings.ToLower(typeName) {
	case "bool":
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return Result{}, fmt.Errorf("cast: key %q value %q is not a valid bool", key, raw)
		}
		return Result{Key: key, Raw: raw, TypeName: "bool", Value: b}, nil
	case "int":
		i, err := strconv.Atoi(raw)
		if err != nil {
			return Result{}, fmt.Errorf("cast: key %q value %q is not a valid int", key, raw)
		}
		return Result{Key: key, Raw: raw, TypeName: "int", Value: i}, nil
	case "float":
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return Result{}, fmt.Errorf("cast: key %q value %q is not a valid float", key, raw)
		}
		return Result{Key: key, Raw: raw, TypeName: "float", Value: f}, nil
	case "string":
		return Result{Key: key, Raw: raw, TypeName: "string", Value: raw}, nil
	default:
		return Result{}, fmt.Errorf("cast: unknown type hint %q for key %q", typeName, key)
	}
}
