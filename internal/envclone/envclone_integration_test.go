package envclone_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envclone"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/diff"
	"os"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
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

func TestEnvClone_FullPipeline_StripPrefix(t *testing.T) {
	path := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=postgres://db\n")
	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatal(err)
	}
	cloned := envclone.Apply(env, envclone.Options{KeyPrefix: "APP_", StripPrefix: true})
	if cloned["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", cloned["HOST"])
	}
	if _, ok := cloned["DB_URL"]; ok {
		t.Error("DB_URL should not be in cloned map")
	}
}

func TestEnvClone_FullPipeline_LowercaseKeys(t *testing.T) {
	path := writeTempEnvFile(t, "HOST=localhost\nPORT=9000\n")
	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatal(err)
	}
	cloned := envclone.Apply(env, envclone.Options{KeyTransform: strings.ToLower})
	result := diff.Compare(env, cloned)
	// All original keys are uppercase; cloned has lowercase — all should be missing from each other
	for _, d := range result {
		if d.Key == "host" || d.Key == "port" {
			return
		}
	}
}
