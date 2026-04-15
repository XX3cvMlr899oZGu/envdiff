package merge_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/merge"
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

func TestMerge_FullPipeline_NoConflicts(t *testing.T) {
	pathA := writeTempEnvFile(t, "APP_NAME=myapp\nDEBUG=false\n")
	pathB := writeTempEnvFile(t, "LOG_LEVEL=info\nPORT=8080\n")

	envA, err := loader.LoadFile(pathA, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile A: %v", err)
	}
	envB, err := loader.LoadFile(pathB, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile B: %v", err)
	}

	res, err := merge.Merge(merge.StrategyFirst, envA, envB)
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if len(res.Env) != 4 {
		t.Errorf("expected 4 keys, got %d", len(res.Env))
	}
}

func TestMerge_FullPipeline_WithConflict(t *testing.T) {
	pathA := writeTempEnvFile(t, "APP_NAME=myapp\nDEBUG=false\n")
	pathB := writeTempEnvFile(t, "APP_NAME=otherapp\nPORT=8080\n")

	envA, err := loader.LoadFile(pathA, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile A: %v", err)
	}
	envB, err := loader.LoadFile(pathB, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile B: %v", err)
	}

	res, err := merge.Merge(merge.StrategyLast, envA, envB)
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if len(res.Conflicts) != 1] != "APP_NAME" {
		t.Errorf("expected conflict on APP_NAME, got %v", res.Conflicts)
	}
	if res.Env["APP_NAME"] != "otherapp" {
		t.Errorf("expected last-wins value 'otherapp', got %q", res.Env["APP_NAME"])
	}
}

func TestMerge_FullPipeline_StrategyError(t *testing.T) {
	pathA := writeTempEnvFile(t, "SECRET=abc\n")
	pathB := writeTempEnvFile(t, "SECRET=xyz\n")

	envA, err := loader.LoadFile(pathA, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile A: %v", err)
	}
	envB, err := loader.LoadFile(pathB, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile B: %v", err)
	}

	_, err = merge.Merge(merge.StrategyError, envA, envB)
	if err == nil {
		t.Fatal("expected error due to conflict, got nil")
	}
}
