package envgroup_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/envgroup"
	"github.com/user/envdiff/internal/loader"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envgroup-*.env")
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

func TestEnvGroup_FullPipeline(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=envdiff\nAPP_DEBUG=true\n")

	env, err := loader.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	groups := envgroup.Apply(env, envgroup.DefaultOptions())
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}

	groupMap := make(map[string]envgroup.Group)
	for _, g := range groups {
		groupMap[g.Name] = g
	}

	if _, ok := groupMap["DB"]; !ok {
		t.Error("expected DB group")
	}
	if _, ok := groupMap["APP"]; !ok {
		t.Error("expected APP group")
	}
	if groupMap["DB"].Keys["DB_HOST"] != "localhost" {
		t.Errorf("unexpected DB_HOST value")
	}
}
