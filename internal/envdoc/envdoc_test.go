package envdoc_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/envdoc"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdoc-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParse_PrecedingComment(t *testing.T) {
	path := writeTempEnv(t, "# Database host\nDB_HOST=localhost\n")
	entries, err := envdoc.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Comment != "Database host" {
		t.Errorf("unexpected comment: %q", entries[0].Comment)
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # HTTP port\n")
	entries, err := envdoc.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].Comment != "HTTP port" {
		t.Errorf("unexpected comment: %q", entries[0].Comment)
	}
	if entries[0].Value != "8080" {
		t.Errorf("unexpected value: %q", entries[0].Value)
	}
}

func TestParse_NoComment(t *testing.T) {
	path := writeTempEnv(t, "SECRET=abc\n")
	entries, err := envdoc.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].Comment != "" {
		t.Errorf("expected empty comment, got %q", entries[0].Comment)
	}
}

func TestToMap_Keys(t *testing.T) {
	entries := []envdoc.Entry{
		{Key: "A", Value: "1", Comment: "alpha"},
		{Key: "B", Value: "2", Comment: ""},
	}
	m := envdoc.ToMap(entries)
	if m["A"] != "alpha" {
		t.Errorf("expected alpha, got %q", m["A"])
	}
	if _, ok := m["B"]; !ok {
		t.Error("expected B to exist in map")
	}
}

func TestParse_BlankLineResetsComment(t *testing.T) {
	path := writeTempEnv(t, "# stale comment\n\nKEY=val\n")
	entries, err := envdoc.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].Comment != "" {
		t.Errorf("expected blank comment after blank line, got %q", entries[0].Comment)
	}
}
