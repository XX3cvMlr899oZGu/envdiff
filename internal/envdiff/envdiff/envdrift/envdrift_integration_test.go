package envdrift_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envdrift"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/snapshot"
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
	return f.Name()
}

func TestEnvDrift_FullPipeline_NoDrift(t *testing.T) {
	path := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=8080\n")

	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	snapshotPath := writeTempEnvFile(t, "")
	if err := snapshot.Save(snapshotPath, env); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := snapshot.Load(snapshotPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	entries := envdrift.Detect(snap, env, envdrift.DefaultOptions())
	if envdrift.HasDrift(entries) {
		t.Errorf("expected no drift, got: %s", envdrift.FormatText(entries))
	}
}

func TestEnvDrift_FullPipeline_DetectsDrift(t *testing.T) {
	snapshotPath := writeTempEnvFile(t, "")
	snap := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "OLD_KEY": "gone"}
	if err := snapshot.Save(snapshotPath, snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	livePath := writeTempEnvFile(t, "APP_HOST=localhost\nAPP_PORT=9090\nNEW_KEY=present\n")
	live, err := loader.LoadFile(livePath, loader.Options{})
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	loaded, err := snapshot.Load(snapshotPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	entries := envdrift.Detect(loaded, live, envdrift.DefaultOptions())
	if !envdrift.HasDrift(entries) {
		t.Fatal("expected drift to be detected")
	}

	statuses := map[string]envdrift.Status{}
	for _, e := range entries {
		statuses[e.Key] = e.Status
	}
	if statuses["APP_PORT"] != envdrift.StatusChanged {
		t.Errorf("APP_PORT should be changed, got %s", statuses["APP_PORT"])
	}
	if statuses["OLD_KEY"] != envdrift.StatusRemoved {
		t.Errorf("OLD_KEY should be removed, got %s", statuses["OLD_KEY"])
	}
	if statuses["NEW_KEY"] != envdrift.StatusAdded {
		t.Errorf("NEW_KEY should be added, got %s", statuses["NEW_KEY"])
	}
}
