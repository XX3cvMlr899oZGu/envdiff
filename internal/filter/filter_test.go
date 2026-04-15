package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/filter"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.local",
		"DB_PASSWORD": "secret",
		"DEBUG":       "true",
	}
}

func TestApplyToMap_NoFilter(t *testing.T) {
	result, err := filter.ApplyToMap(baseEnv(), filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(baseEnv()) {
		t.Errorf("expected %d keys, got %d", len(baseEnv()), len(result))
	}
}

func TestApplyToMap_PrefixFilter(t *testing.T) {
	result, err := filter.ApplyToMap(baseEnv(), filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestApplyToMap_ExcludeKeys(t *testing.T) {
	result, err := filter.ApplyToMap(baseEnv(), filter.Options{Exclude: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
}

func TestApplyToMap_RegexFilter(t *testing.T) {
	result, err := filter.ApplyToMap(baseEnv(), filter.Options{KeyRegex: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApplyToMap_InvalidRegex(t *testing.T) {
	_, err := filter.ApplyToMap(baseEnv(), filter.Options{KeyRegex: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestApplyToMap_PrefixAndExcludeCombined(t *testing.T) {
	result, err := filter.ApplyToMap(baseEnv(), filter.Options{
		Prefix:  "DB_",
		Exclude: []string{"DB_PASSWORD"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}
