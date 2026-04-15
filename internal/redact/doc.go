// Package redact provides utilities for masking sensitive values in
// environment variable maps before they are displayed or exported.
//
// Sensitive keys are identified by matching their names against a set of
// configurable regular expressions. Matched values are replaced with a
// placeholder string (default: "***REDACTED***").
//
// Example usage:
//
//	env := map[string]string{
//		"DB_PASSWORD": "hunter2",
//		"APP_ENV":     "production",
//	}
//
//	redacted, err := redact.Apply(env, redact.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	// redacted["DB_PASSWORD"] == "***REDACTED***"
//	// redacted["APP_ENV"]     == "production"
package redact
