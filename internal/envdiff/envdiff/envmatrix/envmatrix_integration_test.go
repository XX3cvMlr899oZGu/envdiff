package envmatrix_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envdiff/internal/envdiff/envdiff/envmatrix"
	"github.com/yourusername/envdiff/internal/loader"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestEnvMatrix_FullPipeline_MismatchedValues(t *testing.T) {
	dev := writeTempEnvFile(t, "APP=myapp\nDEBUG=true\nPORT=8080\n")
	prod := writeTempEnvFile(t, "APP=myapp\nDEBUG=false\nPORT=443\n")

	devMap, err := loader.LoadFile(dev, loader.Options{})
	if err != nil {
		t.Fatalf("load dev: %v", err)
	}
	prodMap, err := loader.LoadFile(prod, loader.Options{})
	if err != nil {
		t.Fatalf("load prod: %v", err)
	}

	envs := map[string]map[string]string{
		"dev":  devMap,
		"prod": prodMap,
	}
	rows := envmatrix.Build(envs, envmatrix.DefaultOptions())
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	out := envmatrix.FormatText(rows, []string{"dev", "prod"})
	if !strings.Contains(out, "DIFF") {
		t.Error("expected DIFF in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK in output for APP")
	}
}

func TestEnvMatrix_FullPipeline_MissingKey(t *testing.T) {
	dev := writeTempEnvFile(t, "APP=myapp\nDEBUG=true\n")
	prod := writeTempEnvFile(t, "APP=myapp\nDEBUG=true\nSENTRY=https://x\n")

	devMap, _ := loader.LoadFile(dev, loader.Options{})
	prodMap, _ := loader.LoadFile(prod, loader.Options{})

	envs := map[string]map[string]string{
		"dev":  devMap,
		"prod": prodMap,
	}
	rows := envmatrix.Build(envs, envmatrix.DefaultOptions())
	var sentryRow *envmatrix.Row
	for i, r := range rows {
		if r.Key == "SENTRY" {
			sentryRow = &rows[i]
			break
		}
	}
	if sentryRow == nil {
		t.Fatal("expected SENTRY row")
	}
	if len(sentryRow.Missing) != 1 || sentryRow.Missing[0] != "dev" {
		t.Errorf("expected SENTRY missing in dev, got %v", sentryRow.Missing)
	}
}
