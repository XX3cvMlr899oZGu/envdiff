package envindex_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envindex"
	"github.com/user/envdiff/internal/loader"
)

func writeTempEnvFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestEnvIndex_FullPipeline_SharedValues(t *testing.T) {
	prod := writeTempEnvFile(t, ".env.prod",
		"DB_HOST=db.internal\nDB_PORT=5432\nAPP_ENV=production\n")
	staging := writeTempEnvFile(t, ".env.staging",
		"DB_HOST=db.internal\nDB_PORT=5432\nAPP_ENV=staging\n")

	prodEnv, err := loader.LoadFile(prod, nil)
	if err != nil {
		t.Fatalf("load prod: %v", err)
	}
	stagingEnv, err := loader.LoadFile(staging, nil)
	if err != nil {
		t.Fatalf("load staging: %v", err)
	}

	envs := map[string]map[string]string{
		"prod":    prodEnv,
		"staging": stagingEnv,
	}

	opts := envindex.DefaultOptions()
	opts.ExcludeUnique = true
	entries := envindex.Build(envs, opts)

	if len(entries) == 0 {
		t.Fatal("expected shared entries")
	}

	for _, e := range entries {
		if len(e.Occurrences) < 2 {
			t.Errorf("ExcludeUnique=true but entry %q has only 1 env", e.Value)
		}
	}

	out := envindex.FormatText(entries)
	if !strings.Contains(out, "db.internal") {
		t.Errorf("expected 'db.internal' in formatted output:\n%s", out)
	}
	if !strings.Contains(out, "5432") {
		t.Errorf("expected '5432' in formatted output:\n%s", out)
	}
}

func TestEnvIndex_FullPipeline_AllUnique(t *testing.T) {
	f := writeTempEnvFile(t, ".env", "KEY1=alpha\nKEY2=beta\n")
	env, err := loader.LoadFile(f, nil)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	envs := map[string]map[string]string{"only": env}
	opts := envindex.DefaultOptions()
	opts.ExcludeUnique = true
	entries := envindex.Build(envs, opts)
	if len(entries) != 0 {
		t.Errorf("expected empty result when all values are unique, got %d", len(entries))
	}
}
