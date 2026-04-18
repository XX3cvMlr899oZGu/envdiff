package envpatch_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/envpatch"
	"github.com/yourorg/envdiff/internal/loader"
	"github.com/yourorg/envdiff/internal/diff"
	"os"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envpatch-*.env")
	if err != nil { t.Fatal(err) }
	if _, err := f.WriteString(content); err != nil { t.Fatal(err) }
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestEnvPatch_FullPipeline_SetAndUnset(t *testing.T) {
	path := writeTempEnvFile(t, "APP_ENV=production\nDB_HOST=localhost\nSECRET=abc\n")
	env, err := loader.LoadFile(path, nil)
	if err != nil { t.Fatal(err) }

	ops := []envpatch.Op{
		{Type: envpatch.OpSet, Key: "APP_ENV", Value: "staging"},
		{Type: envpatch.OpUnset, Key: "SECRET"},
	}
	patched, _, err := envpatch.Apply(env, ops)
	if err != nil { t.Fatal(err) }

	if patched["APP_ENV"] != "staging" { t.Errorf("expected staging") }
	if _, ok := patched["SECRET"]; ok { t.Error("SECRET should be removed") }
}

func TestEnvPatch_FullPipeline_RenameAndDiff(t *testing.T) {
	base := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	env, err := loader.LoadFile(base, nil)
	if err != nil { t.Fatal(err) }

	ops := []envpatch.Op{
		{Type: envpatch.OpRename, Key: "DB_HOST", To: "DATABASE_HOST"},
	}
	patched, _, err := envpatch.Apply(env, ops)
	if err != nil { t.Fatal(err) }

	results := diff.Compare(env, patched)
	found := false
	for _, r := range results {
		if r.Key == "DATABASE_HOST" { found = true }
	}
	if !found { t.Error("expected DATABASE_HOST in diff results") }
}
