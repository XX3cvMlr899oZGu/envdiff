package envsort_test

import (
	"testing"

	"github.com/user/envdiff/internal/envsort"
	"github.com/user/envdiff/internal/loader"
	"os"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestEnvSort_FullPipeline_Ascending(t *testing.T) {
	path := writeTempEnvFile(t, "ZEBRA=1\nAPPLE=2\nMIDDLE=3\n")
	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatal(err)
	}
	entries := envsort.Apply(env, envsort.DefaultOptions())
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "APPLE" {
		t.Errorf("first key should be APPLE, got %s", entries[0].Key)
	}
	if entries[2].Key != "ZEBRA" {
		t.Errorf("last key should be ZEBRA, got %s", entries[2].Key)
	}
}

func TestEnvSort_FullPipeline_PrefixPriority(t *testing.T) {
	path := writeTempEnvFile(t, "APP_NAME=foo\nDB_HOST=localhost\nDB_PASS=secret\nZONE=us\n")
	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatal(err)
	}
	opts := envsort.Options{Order: envsort.Ascending, Prefix: "DB_"}
	entries := envsort.Apply(env, opts)
	if entries[0].Key != "DB_HOST" && entries[0].Key != "DB_PASS" {
		t.Errorf("expected DB_ key first, got %s", entries[0].Key)
	}
}
