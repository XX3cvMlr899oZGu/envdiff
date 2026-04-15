package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

var sampleResults = []diff.Result{
	{Key: "APP_ENV", Status: diff.StatusEqual, ValueA: "prod", ValueB: "prod"},
	{Key: "DB_HOST", Status: diff.StatusMismatch, ValueA: "localhost", ValueB: "db.prod"},
	{Key: "SECRET", Status: diff.StatusMissingInB, ValueA: "abc", ValueB: ""},
	{Key: "NEW_KEY", Status: diff.StatusMissingInA, ValueA: "", ValueB: "xyz"},
}

func TestWriteText_ContainsDiffLines(t *testing.T) {
	var buf bytes.Buffer
	err := report.Write(&buf, sampleResults, ".env.dev", ".env.prod", report.Options{Format: report.FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected SECRET in output")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("expected NEW_KEY in output")
	}
}

func TestWriteText_HidesEqualByDefault(t *testing.T) {
	var buf bytes.Buffer
	_ = report.Write(&buf, sampleResults, "a", "b", report.Options{Format: report.FormatText})
	if strings.Contains(buf.String(), "[=]") {
		t.Error("equal keys should be hidden by default")
	}
}

func TestWriteText_ShowEqual(t *testing.T) {
	var buf bytes.Buffer
	_ = report.Write(&buf, sampleResults, "a", "b", report.Options{Format: report.FormatText, ShowEqual: true})
	if !strings.Contains(buf.String(), "[=]") {
		t.Error("expected equal keys when ShowEqual is true")
	}
}

func TestWriteJSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	err := report.Write(&buf, sampleResults, "a", "b", report.Options{Format: report.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["file_a"] != "a" || out["file_b"] != "b" {
		t.Error("expected file_a and file_b fields")
	}
	results, ok := out["results"].([]interface{})
	if !ok || len(results) != 4 {
		t.Errorf("expected 4 results, got %v", len(results))
	}
}
