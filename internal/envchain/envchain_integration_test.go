package envchain_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/envchain"
	"github.com/yourorg/envdiff/internal/loader"
	"github.com/yourorg/envdiff/internal/filter"
	"github.com/yourorg/envdiff/internal/envtrim"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestEnvChain_FullPipeline_FilterAndTrim(t *testing.T) {
	path := writeTempEnvFile(t, `
APP_HOST=  localhost  
APP_PORT=  8080  
DB_HOST=  db.internal  
DEBUG=true
`)

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	chain := envchain.New().
		Add("filter-prefix", func(e map[string]string) (map[string]string, error) {
			return filter.ApplyToMap(e, filter.Options{Prefix: "APP_"}), nil
		}).
		Add("trim", func(e map[string]string) (map[string]string, error) {
			return envtrim.Apply(e, envtrim.DefaultOptions()), nil
		})

	out, err := chain.Run(env)
	if err != nil {
		t.Fatalf("chain run: %v", err)
	}

	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should have been filtered out")
	}
	if _, ok := out["DEBUG"]; ok {
		t.Error("DEBUG should have been filtered out")
	}
	if strings.TrimSpace(out["APP_HOST"]) != out["APP_HOST"] {
		t.Errorf("APP_HOST value should be trimmed, got %q", out["APP_HOST"])
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", out["APP_HOST"])
	}
}

func TestEnvChain_FullPipeline_EmptyInput(t *testing.T) {
	chain := envchain.New().
		Add("trim", func(e map[string]string) (map[string]string, error) {
			return envtrim.Apply(e, envtrim.DefaultOptions()), nil
		})

	out, err := chain.Run(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}
