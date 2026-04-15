package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/export"
)

var sampleResults = []diff.Result{
	{Key: "DB_HOST", Status: diff.StatusEqual, ValueA: "localhost", ValueB: "localhost"},
	{Key: "API_KEY", Status: diff.StatusMissing, ValueA: "abc123", ValueB: ""},
	{Key: "PORT", Status: diff.StatusMismatch, ValueA: "8080", ValueB: "9090"},
}

func TestWriteCSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, sampleResults, export.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "key,status,value_a,value_b") {
		t.Errorf("expected CSV header, got:\n%s", output)
	}
}

func TestWriteCSV_ContainsRows(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, sampleResults, export.FormatCSV)
	output := buf.String()
	if !strings.Contains(output, "API_KEY") {
		t.Errorf("expected API_KEY in output, got:\n%s", output)
	}
	if !strings.Contains(output, "PORT") {
		t.Errorf("expected PORT in output, got:\n%s", output)
	}
}

func TestWriteMarkdown_ContainsTable(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, sampleResults, export.FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "| Key | Status |") {
		t.Errorf("expected markdown header, got:\n%s", output)
	}
	if !strings.Contains(output, "DB_HOST") {
		t.Errorf("expected DB_HOST row, got:\n%s", output)
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, sampleResults, export.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestWriteCSV_SortedByKey(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, sampleResults, export.FormatCSV)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// lines[0] is header; keys should be sorted: API_KEY, DB_HOST, PORT
	if !strings.HasPrefix(lines[1], "API_KEY") {
		t.Errorf("expected first data row to be API_KEY, got: %s", lines[1])
	}
}
