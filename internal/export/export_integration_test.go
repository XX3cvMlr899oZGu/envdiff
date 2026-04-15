package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/export"
)

func TestExportCSV_FullPipeline(t *testing.T) {
	envA := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
		"SECRET": "s3cr3t",
	}
	envB := map[string]string{
		"HOST": "localhost",
		"PORT": "9090",
	}

	results := diff.Compare(envA, envB)

	var buf bytes.Buffer
	if err := export.Write(&buf, results, export.FormatCSV); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "PORT") {
		t.Error("expected PORT in CSV output")
	}
	if !strings.Contains(output, "SECRET") {
		t.Error("expected SECRET in CSV output")
	}
	if !strings.Contains(output, string(diff.StatusMismatch)) {
		t.Errorf("expected mismatch status in output, got:\n%s", output)
	}
}

func TestExportMarkdown_FullPipeline(t *testing.T) {
	envA := map[string]string{"ALPHA": "1", "BETA": "2"}
	envB := map[string]string{"ALPHA": "1"}

	results := diff.Compare(envA, envB)

	var buf bytes.Buffer
	if err := export.Write(&buf, results, export.FormatMarkdown); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "BETA") {
		t.Error("expected BETA in markdown output")
	}
	if !strings.Contains(output, "|") {
		t.Error("expected markdown table pipes in output")
	}
}
