package redact_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/redact"
)

func TestRedactThenCompare(t *testing.T) {
	prod := map[string]string{
		"DB_PASSWORD": "prod-secret",
		"APP_NAME":    "myapp",
		"PORT":        "8080",
	}
	staging := map[string]string{
		"DB_PASSWORD": "staging-secret",
		"APP_NAME":    "myapp",
		"PORT":        "9090",
	}

	opts := redact.DefaultOptions()

	redactedProd, err := redact.Apply(prod, opts)
	if err != nil {
		t.Fatalf("redact prod: %v", err)
	}
	redactedStaging, err := redact.Apply(staging, opts)
	if err != nil {
		t.Fatalf("redact staging: %v", err)
	}

	// After redaction, DB_PASSWORD should appear equal (both redacted)
	results := diff.Compare(redactedProd, redactedStaging)

	for _, r := range results {
		if r.Key == "DB_PASSWORD" && r.Status != "equal" {
			t.Errorf("expected DB_PASSWORD to appear equal after redaction, got status %q", r.Status)
		}
		if r.Key == "PORT" && r.Status != "mismatch" {
			t.Errorf("expected PORT to be mismatch, got status %q", r.Status)
		}
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	original := map[string]string{
		"DB_PASSWORD": "real-secret",
		"HOST":        "localhost",
	}

	_, err := redact.Apply(original, redact.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if original["DB_PASSWORD"] != "real-secret" {
		t.Error("Apply must not mutate the original map")
	}
}
