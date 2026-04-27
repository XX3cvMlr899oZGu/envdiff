package envprofile

import (
	"testing"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"REDIS_HOST":  "redis",
	"APP_DEBUG":   "true",
	"APP_VERSION": "1.0.0",
	"SECRET_KEY":  "abc123",
}

func TestApply_ExactMatch(t *testing.T) {
	p := Profile{Name: "db", Patterns: []string{"DB_HOST", "DB_PORT"}}
	out, err := Apply(sampleEnv, p, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestApply_PrefixGlob(t *testing.T) {
	p := Profile{Name: "app", Patterns: []string{"APP_*"}}
	out, err := Apply(sampleEnv, p, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_IncludeUnmatched(t *testing.T) {
	p := Profile{Name: "db", Patterns: []string{"DB_*"}}
	opts := Options{IncludeUnmatched: true}
	out, err := Apply(sampleEnv, p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(sampleEnv) {
		t.Errorf("expected all %d keys, got %d", len(sampleEnv), len(out))
	}
}

func TestApply_EmptyProfileNameReturnsError(t *testing.T) {
	p := Profile{Name: "", Patterns: []string{"DB_*"}}
	_, err := Apply(sampleEnv, p, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty profile name")
	}
}

func TestApply_NoMatchReturnsEmptyMap(t *testing.T) {
	p := Profile{Name: "none", Patterns: []string{"UNKNOWN_*"}}
	out, err := Apply(sampleEnv, p, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}

func TestMatchedKeys_Sorted(t *testing.T) {
	p := Profile{Name: "db", Patterns: []string{"DB_*"}}
	keys := MatchedKeys(sampleEnv, p)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PORT" {
		t.Errorf("unexpected order: %v", keys)
	}
}
