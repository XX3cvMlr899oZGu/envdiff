package envdefault_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/envdefault"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/diff"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdefault-*.env")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestEnvDefault_FullPipeline_FillsMissing(t *testing.T) {
	envFile := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	defsFile := writeTempEnvFile(t, "APP_HOST=0.0.0.0\nAPP_PORT=80\nAPP_DEBUG=false\n")

	env, err := loader.LoadFile(envFile, loader.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defs, err := loader.LoadFile(defsFile, loader.Options{})
	if err != nil {
		t.Fatal(err)
	}

	result := envdefault.Apply(env, defs, envdefault.DefaultOptions())

	if result["APP_HOST"] != "localhost" {
		t.Errorf("APP_HOST should not be overwritten, got %s", result["APP_HOST"])
	}
	if result["APP_DEBUG"] != "false" {
		t.Errorf("APP_DEBUG should be filled from defaults, got %s", result["APP_DEBUG"])
	}

	missing := envdefault.MissingKeys(env, defs)
	if len(missing) != 1 || missing[0] != "APP_DEBUG" {
		t.Errorf("expected [APP_DEBUG] missing, got %v", missing)
	}

	results := diff.Compare(result, defs)
	for _, r := range results {
		if r.Key == "APP_DEBUG" && r.Status != "equal" {
			t.Errorf("APP_DEBUG should be equal after apply")
		}
	}
}
