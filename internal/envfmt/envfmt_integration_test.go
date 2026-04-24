package envfmt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/envfmt"
	"github.com/yourorg/envdiff/internal/parser"
	"os"
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

func TestEnvFmt_FullPipeline_RoundTrip(t *testing.T) {
	path := writeTempEnvFile(t, "ZEBRA=z\nALPHA=a\nMIDDLE=m\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	opts := envfmt.DefaultOptions()
	if err := envfmt.Apply(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected sorted first key ALPHA, got %s", lines[0])
	}
}

func TestEnvFmt_FullPipeline_QuotedOutput(t *testing.T) {
	path := writeTempEnvFile(t, "HOST=localhost\nPORT=5432\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	opts := envfmt.DefaultOptions()
	opts.QuoteStyle = envfmt.QuoteDouble
	if err := envfmt.Apply(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `="`) {
		t.Errorf("expected double-quoted values in output:\n%s", out)
	}
}
