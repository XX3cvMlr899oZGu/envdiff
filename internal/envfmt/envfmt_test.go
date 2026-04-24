package envfmt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/envfmt"
)

func TestApply_SortedOutput(t *testing.T) {
	env := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MIDDLE": "m"}
	var buf bytes.Buffer
	opts := envfmt.DefaultOptions()
	if err := envfmt.Apply(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first line to start with ALPHA, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line to start with ZEBRA, got %s", lines[2])
	}
}

func TestApply_DoubleQuotes(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	var buf bytes.Buffer
	opts := envfmt.DefaultOptions()
	opts.QuoteStyle = envfmt.QuoteDouble
	if err := envfmt.Apply(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if got != `KEY="value"` {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestApply_SingleQuotes(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	var buf bytes.Buffer
	opts := envfmt.DefaultOptions()
	opts.QuoteStyle = envfmt.QuoteSingle
	if err := envfmt.Apply(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if got != `KEY='value'` {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestApply_InlineComments(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	var buf bytes.Buffer
	opts := envfmt.DefaultOptions()
	opts.Comments = map[string]string{"PORT": "HTTP port"}
	if err := envfmt.Apply(&buf, env, opts); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "# HTTP port") {
		t.Errorf("expected inline comment, got: %q", got)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	var buf bytes.Buffer
	if err := envfmt.Apply(&buf, map[string]string{}, envfmt.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}
