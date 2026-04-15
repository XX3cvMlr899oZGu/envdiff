package rename_test

import (
	"testing"

	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/rename"
	"os"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-rename-*.env")
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

func TestRename_FullPipeline_PrefixSwap(t *testing.T) {
	path := writeTempEnvFile(t, "DEV_DB=postgres\nDEV_HOST=localhost\nAPP_NAME=myapp\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	opts := rename.Options{
		Map:       make(map[string]string),
		OldPrefix: "DEV_",
		NewPrefix: "PROD_",
	}
	res, err := rename.Apply(env, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if res.Env["PROD_DB"] != "postgres" {
		t.Errorf("expected PROD_DB=postgres, got %q", res.Env["PROD_DB"])
	}
	if res.Env["PROD_HOST"] != "localhost" {
		t.Errorf("expected PROD_HOST=localhost")
	}
	if res.Env["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged")
	}
	if len(res.Changes) != 2 {
		t.Errorf("expected 2 renames, got %d", len(res.Changes))
	}
}

func TestRename_FullPipeline_ExplicitMap(t *testing.T) {
	path := writeTempEnvFile(t, "OLD_TOKEN=abc123\nKEEP=yes\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	opts := rename.DefaultOptions()
	opts.Map["OLD_TOKEN"] = "API_TOKEN"
	res, err := rename.Apply(env, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if res.Env["API_TOKEN"] != "abc123" {
		t.Errorf("expected API_TOKEN=abc123")
	}
	if _, ok := res.Env["OLD_TOKEN"]; ok {
		t.Error("OLD_TOKEN should be gone")
	}
}
