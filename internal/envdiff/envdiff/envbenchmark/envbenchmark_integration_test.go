package envbenchmark_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envdiff/envbenchmark"
	"github.com/user/envdiff/internal/loader"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestEnvBenchmark_FullPipeline_Noop(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 20; i++ {
		sb.WriteString(fmt.Sprintf("KEY_%d=value_%d\n", i, i))
	}
	p := writeTempEnvFile(t, sb.String())

	env, err := loader.LoadFile(p, nil)
	if err != nil {
		t.Fatalf("load file: %v", err)
	}

	opts := envbenchmark.Options{Runs: 3}
	r, err := envbenchmark.Run("load-noop", env, func(m map[string]string) error { return nil }, opts)
	if err != nil {
		t.Fatalf("benchmark run: %v", err)
	}
	if r.KeyCount != 20 {
		t.Errorf("expected 20 keys, got %d", r.KeyCount)
	}
	if r.Label != "load-noop" {
		t.Errorf("unexpected label: %q", r.Label)
	}
}

func TestEnvBenchmark_FullPipeline_FormatOutput(t *testing.T) {
	p := writeTempEnvFile(t, "A=1\nB=2\nC=3\n")
	env, err := loader.LoadFile(p, nil)
	if err != nil {
		t.Fatalf("load file: %v", err)
	}

	r, _ := envbenchmark.Run("fmt-test", env, func(m map[string]string) error { return nil }, envbenchmark.DefaultOptions())
	out := envbenchmark.FormatText([]envbenchmark.Result{r})

	if !strings.Contains(out, "fmt-test") {
		t.Errorf("expected label in formatted output, got:\n%s", out)
	}
	if !strings.Contains(out, "keys") {
		t.Errorf("expected header in formatted output, got:\n%s", out)
	}
}
