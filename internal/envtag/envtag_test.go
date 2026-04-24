package envtag_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/envtag"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"REDIS_HOST":  "redis",
	"APP_SECRET":  "s3cr3t",
	"APP_VERSION": "1.0",
}

func TestApply_NoTags(t *testing.T) {
	opts := envtag.DefaultOptions()
	res, err := envtag.Apply(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["_untagged"]) != len(baseEnv) {
		t.Errorf("expected all keys untagged, got %d", len(res["_untagged"]))
	}
}

func TestApply_GroupsByTag(t *testing.T) {
	opts := envtag.Options{
		Tags: []envtag.Tag{
			{Name: "database", Keys: []string{"DB_HOST", "DB_PORT"}},
			{Name: "cache", Keys: []string{"REDIS_HOST"}},
		},
	}
	res, err := envtag.Apply(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["database"]) != 2 {
		t.Errorf("expected 2 database keys, got %d", len(res["database"]))
	}
	if len(res["cache"]) != 1 {
		t.Errorf("expected 1 cache key, got %d", len(res["cache"]))
	}
	if len(res["_untagged"]) != 2 {
		t.Errorf("expected 2 untagged keys, got %d", len(res["_untagged"]))
	}
}

func TestApply_KeyInMultipleTags(t *testing.T) {
	opts := envtag.Options{
		Tags: []envtag.Tag{
			{Name: "secrets", Keys: []string{"APP_SECRET"}},
			{Name: "app", Keys: []string{"APP_SECRET", "APP_VERSION"}},
		},
	}
	res, err := envtag.Apply(baseEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res["secrets"]["APP_SECRET"]; !ok {
		t.Error("APP_SECRET should be in secrets group")
	}
	if _, ok := res["app"]["APP_SECRET"]; !ok {
		t.Error("APP_SECRET should also be in app group")
	}
}

func TestApply_EmptyTagNameReturnsError(t *testing.T) {
	opts := envtag.Options{
		Tags: []envtag.Tag{{Name: "", Keys: []string{"DB_HOST"}}},
	}
	_, err := envtag.Apply(baseEnv, opts)
	if err == nil {
		t.Error("expected error for empty tag name")
	}
}

func TestApply_NilEnv(t *testing.T) {
	res, err := envtag.Apply(nil, envtag.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected empty result for nil env")
	}
}

func TestKeysForTag_ReturnsKeys(t *testing.T) {
	res := envtag.Result{
		"database": {"DB_HOST": "localhost", "DB_PORT": "5432"},
	}
	keys := envtag.KeysForTag(res, "database")
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestKeysForTag_MissingTag(t *testing.T) {
	res := envtag.Result{}
	keys := envtag.KeysForTag(res, "missing")
	if keys != nil {
		t.Errorf("expected nil for missing tag")
	}
}
