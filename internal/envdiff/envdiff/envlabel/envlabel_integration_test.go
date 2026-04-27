package envlabel_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envdiff/internal/loader"
	"github.com/yourusername/envdiff/internal/envdiff/envdiff/envlabel"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestEnvLabel_FullPipeline_AddsLabels(t *testing.T) {
	path := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=9000\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	opts := envlabel.DefaultOptions()
	opts.Labels = map[string]string{"env": "staging", "owner": "infra"}

	le, err := envlabel.Apply(env, nil, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if le.Labels["env"] != "staging" {
		t.Errorf("expected label env=staging")
	}
	if le.Env["APP_HOST"] != "localhost" {
		t.Errorf("env key APP_HOST not preserved")
	}

	text := envlabel.FormatText(le)
	if text == "(no labels)" {
		t.Error("expected non-empty label output")
	}
}

func TestEnvLabel_FullPipeline_ConflictBlocked(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=db.local\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	existing := map[string]string{"env": "production"}
	opts := envlabel.DefaultOptions()
	opts.Labels = map[string]string{"env": "staging"}

	_, err = envlabel.Apply(env, existing, opts)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}
