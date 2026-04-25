package envpivot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/envdiff/envdiff/envpivot"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return filepath.Clean(f.Name())
}

func TestEnvPivot_FullPipeline_BasicTable(t *testing.T) {
	devPath := writeTempEnvFile(t, "APP=myapp\nDB=localhost\nDEBUG=true\n")
	prodPath := writeTempEnvFile(t, "APP=myapp\nDB=prod.db\n")

	devEnv, err := loader.LoadFile(devPath, loader.Options{})
	if err != nil {
		t.Fatalf("load dev: %v", err)
	}
	prodEnv, err := loader.LoadFile(prodPath, loader.Options{})
	if err != nil {
		t.Fatalf("load prod: %v", err)
	}

	envs := map[string]map[string]string{
		"dev":  devEnv,
		"prod": prodEnv,
	}

	rows := envpivot.Pivot(envs, envpivot.DefaultOptions())
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}

	for _, r := range rows {
		if r.Key == "APP" && !r.AllEqual() {
			t.Error("APP should be equal across envs")
		}
		if r.Key == "DB" && r.AllEqual() {
			t.Error("DB should differ across envs")
		}
	}
}

func TestEnvPivot_FullPipeline_ExcludeMissing(t *testing.T) {
	devPath := writeTempEnvFile(t, "APP=myapp\nDEBUG=true\n")
	prodPath := writeTempEnvFile(t, "APP=myapp\n")

	devEnv, _ := loader.LoadFile(devPath, loader.Options{})
	prodEnv, _ := loader.LoadFile(prodPath, loader.Options{})

	envs := map[string]map[string]string{
		"dev":  devEnv,
		"prod": prodEnv,
	}

	rows := envpivot.Pivot(envs, envpivot.Options{ExcludeMissing: true})
	for _, r := range rows {
		if r.Key == "DEBUG" {
			t.Error("DEBUG should be excluded when ExcludeMissing=true")
		}
	}
	if len(rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(rows))
	}
}
