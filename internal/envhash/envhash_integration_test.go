package envhash_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/user/envdiff/internal/envhash"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/filter"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := fmt.Fprint(f, content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	_ = f.Close()
	return f.Name()
}

func TestEnvHash_FullPipeline_SameFiles(t *testing.T) {
	content := "APP_HOST=localhost\nAPP_PORT=8080\nDEBUG=false\n"
	path1 := writeTempEnvFile(t, content)
	path2 := writeTempEnvFile(t, content)

	env1, err := loader.LoadFile(path1, filter.ApplyToMap(nil))
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	env2, err := loader.LoadFile(path2, filter.ApplyToMap(nil))
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	eq, err := envhash.Equal(env1, env2, envhash.DefaultOptions())
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if !eq {
		t.Error("expected Equal to be true for identical file contents")
	}
}

func TestEnvHash_FullPipeline_DifferentFiles(t *testing.T) {
	path1 := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	path2 := writeTempEnvFile(t, "APP_HOST=remotehost\nAPP_PORT=9090\n")

	env1, err := loader.LoadFile(path1, filter.ApplyToMap(nil))
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	env2, err := loader.LoadFile(path2, filter.ApplyToMap(nil))
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	eq, err := envhash.Equal(env1, env2, envhash.DefaultOptions())
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if eq {
		t.Error("expected Equal to be false for different file contents")
	}
}
