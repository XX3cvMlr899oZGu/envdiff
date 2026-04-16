// Package envtype infers the likely type of environment variable values.
package envtype

import (
	"regexp"
	"strconv"
	"strings"
)

// Type represents an inferred value type.
type Type string

const (
	TypeBool   Type = "bool"
	TypeInt    Type = "int"
	TypeFloat  Type = "float"
	TypeURL    Type = "url"
	TypeEmail  Type = "email"
	TypePath   Type = "path"
	TypeString Type = "string"
)

var (
	urlRe   = regexp.MustCompile(`(?i)^https?://`)
	emailRe = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	pathRe  = regexp.MustCompile(`^[/~.]`)
)

// Infer returns the most specific Type that matches the given value.
func Infer(value string) Type {
	v := strings.TrimSpace(value)

	if _, err := strconv.ParseBool(v); err == nil {
		return TypeBool
	}
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return TypeInt
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return TypeFloat
	}
	if urlRe.MatchString(v) {
		return TypeURL
	}
	if emailRe.MatchString(v) {
		return TypeEmail
	}
	if pathRe.MatchString(v) {
		return TypePath
	}
	return TypeString
}

// InferAll returns a map of key → inferred Type for every entry.
func InferAll(env map[string]string) map[string]Type {
	result := make(map[string]Type, len(env))
	for k, v := range env {
		result[k] = Infer(v)
	}
	return result
}
