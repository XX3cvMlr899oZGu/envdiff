package cast_test

import (
	"os"
	"testing"

	"github.com/yourorg/envdiff/internal/cast"
	"github.com/yourorg/envdiff/internal/loader"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestCast_FullPipeline(t *testing.T) {
	path := writeTempEnvFile(t, "PORT=9090\nDEBUG=false\nAPP_NAME=envdiff\nRATIO=2.71\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("loader error: %v", err)
	}

	opts := cast.DefaultOptions()
	opts.TypeHints = map[string]string{
		"PORT":  "int",
		"DEBUG": "bool",
		"RATIO": "float",
	}

	results, err := cast.Apply(env, opts)
	if err != nil {
		t.Fatalf("cast error: %v", err)
	}

	byKey := map[string]cast.Result{}
	for _, r := range results {
		byKey[r.Key] = r
	}

	if byKey["PORT"].Value != 9090 {
		t.Errorf("PORT: expected 9090, got %v", byKey["PORT"].Value)
	}
	if byKey["DEBUG"].Value != false {
		t.Errorf("DEBUG: expected false, got %v", byKey["DEBUG"].Value)
	}
	if byKey["APP_NAME"].TypeName != "string" {
		t.Errorf("APP_NAME: expected string type, got %v", byKey["APP_NAME"].TypeName)
	}
}
