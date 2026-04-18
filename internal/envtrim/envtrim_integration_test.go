package envtrim_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/envtrim"
	"github.com/user/envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envtrim-*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestEnvTrim_FullPipeline_TrimAndUppercase(t *testing.T) {
	path := writeTempEnvFile(t, "db_host=  localhost  \napp_port=  9000\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	opts := envtrim.Options{TrimValues: true, TrimKeys: true, UppercaseKeys: true}
	out := envtrim.Apply(env, opts)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["DB_HOST"])
	}
	if out["APP_PORT"] != "9000" {
		t.Errorf("expected '9000', got %q", out["APP_PORT"])
	}
}

func TestEnvTrim_FullPipeline_NoOp(t *testing.T) {
	path := writeTempEnvFile(t, "KEY=value\nOTHER=123\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	out := envtrim.Apply(env, envtrim.DefaultOptions())
	if out["KEY"] != "value" {
		t.Errorf("unexpected value: %q", out["KEY"])
	}
	if len(out) != len(env) {
		t.Errorf("length mismatch: got %d want %d", len(out), len(env))
	}
}
