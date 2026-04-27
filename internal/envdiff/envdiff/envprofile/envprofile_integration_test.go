package envprofile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/envdiff/envdiff/envprofile"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestEnvProfile_FullPipeline_PrefixGlob(t *testing.T) {
	path := writeTempEnvFile(t, `
DB_HOST=localhost
DB_PORT=5432
APP_DEBUG=true
APP_VERSION=2.0.0
SECRET_KEY=s3cr3t
`)
	env, err := loader.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	p := envprofile.Profile{Name: "database", Patterns: []string{"DB_*"}}
	out, err := envprofile.Apply(env, p, envprofile.DefaultOptions())
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 DB keys, got %d", len(out))
	}
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should not be in database profile")
	}
}

func TestEnvProfile_FullPipeline_MatchedKeysSorted(t *testing.T) {
	path := writeTempEnvFile(t, `
APP_NAME=envdiff
APP_ENV=production
APP_PORT=8080
OTHER=value
`)
	env, err := loader.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	p := envprofile.Profile{Name: "app", Patterns: []string{"APP_*"}}
	keys := envprofile.MatchedKeys(env, p)
	if len(keys) != 3 {
		t.Fatalf("expected 3 APP keys, got %d", len(keys))
	}
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}
