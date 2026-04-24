package envtag_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envdiff/internal/envtag"
	"github.com/yourusername/envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestEnvTag_FullPipeline_GroupsByTag(t *testing.T) {
	path := writeTempEnvFile(t, `
DB_HOST=localhost
DB_PORT=5432
REDIS_HOST=redis
APP_SECRET=topsecret
APP_VERSION=2.0
`)

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	opts := envtag.Options{
		Tags: []envtag.Tag{
			{Name: "database", Keys: []string{"DB_HOST", "DB_PORT"}},
			{Name: "security", Keys: []string{"APP_SECRET"}},
		},
	}

	res, err := envtag.Apply(env, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if res["database"]["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost in database group")
	}
	if res["database"]["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432 in database group")
	}
	if res["security"]["APP_SECRET"] != "topsecret" {
		t.Errorf("expected APP_SECRET in security group")
	}
	if len(res["_untagged"]) != 2 {
		t.Errorf("expected 2 untagged keys (REDIS_HOST, APP_VERSION), got %d", len(res["_untagged"]))
	}
}

func TestEnvTag_FullPipeline_AllUntagged(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\nBAZ=qux\n")

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	res, err := envtag.Apply(env, envtag.DefaultOptions())
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if len(res["_untagged"]) != 2 {
		t.Errorf("expected 2 untagged keys, got %d", len(res["_untagged"]))
	}
}
