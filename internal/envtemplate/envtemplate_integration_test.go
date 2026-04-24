package envtemplate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/envtemplate"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/filter"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestEnvTemplate_FullPipeline_Substitution(t *testing.T) {
	envFile := writeTempEnvFile(t, "DB_USER=admin\nDB_PASS=secret\n")
	tmplFile := writeTempEnvFile(t, "DSN=postgres://{{.DB_USER}}:{{.DB_PASS}}@localhost/app\nAPP=myapp\n")

	envMap, err := loader.LoadFile(envFile, filter.ApplyToMap, nil)
	if err != nil {
		t.Fatalf("LoadFile env: %v", err)
	}
	tmplMap, err := loader.LoadFile(tmplFile, filter.ApplyToMap, nil)
	if err != nil {
		t.Fatalf("LoadFile tmpl: %v", err)
	}

	out, err := envtemplate.Apply(tmplMap, envMap, envtemplate.DefaultOptions())
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	want := "postgres://admin:secret@localhost/app"
	if out["DSN"] != want {
		t.Errorf("DSN: got %q, want %q", out["DSN"], want)
	}
	if out["APP"] != "myapp" {
		t.Errorf("APP: got %q, want %q", out["APP"], "myapp")
	}
}

func TestEnvTemplate_FullPipeline_MissingKeyError(t *testing.T) {
	envFile := writeTempEnvFile(t, "DB_USER=admin\n")
	tmplFile := writeTempEnvFile(t, "DSN=postgres://{{.DB_USER}}:{{.DB_PASS}}@localhost/app\n")

	envMap, err := loader.LoadFile(envFile, filter.ApplyToMap, nil)
	if err != nil {
		t.Fatalf("LoadFile env: %v", err)
	}
	tmplMap, err := loader.LoadFile(tmplFile, filter.ApplyToMap, nil)
	if err != nil {
		t.Fatalf("LoadFile tmpl: %v", err)
	}

	_, err = envtemplate.Apply(tmplMap, envMap, envtemplate.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing DB_PASS, got nil")
	}
}
