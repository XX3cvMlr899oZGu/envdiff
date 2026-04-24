package envscope

import (
	"testing"
)

var base = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PORT":     "5432",
	"LOG_LEVEL":   "info",
}

func TestApply_MatchesPrefix(t *testing.T) {
	scope := Scope{Name: "app", Prefix: "APP_"}
	result, err := Apply(base, scope, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost")
	}
}

func TestApply_StripPrefix(t *testing.T) {
	scope := Scope{Name: "db", Prefix: "DB_"}
	opts := Options{StripPrefix: true}
	result, err := Apply(base, scope, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["HOST"]; !ok {
		t.Errorf("expected key HOST after stripping prefix")
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Errorf("expected DB_HOST to be stripped")
	}
}

func TestApply_EmptyPrefix_ReturnsError(t *testing.T) {
	scope := Scope{Name: "bad", Prefix: ""}
	_, err := Apply(base, scope, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestApply_CaseFold(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "app_port": "9090"}
	scope := Scope{Name: "app", Prefix: "APP_"}
	opts := Options{CaseFold: true}
	result, err := Apply(env, scope, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys with case fold, got %d", len(result))
	}
}

func TestApplyAll_MultipleScopes(t *testing.T) {
	scopes := []Scope{
		{Name: "app", Prefix: "APP_"},
		{Name: "db", Prefix: "DB_"},
	}
	out, err := ApplyAll(base, scopes, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out["app"]) != 2 {
		t.Errorf("expected 2 app keys, got %d", len(out["app"]))
	}
	if len(out["db"]) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(out["db"]))
	}
}

func TestApply_NoMatch_ReturnsEmpty(t *testing.T) {
	scope := Scope{Name: "cache", Prefix: "CACHE_"}
	result, err := Apply(base, scope, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}
