package envmatrix

import (
	"strings"
	"testing"
)

var envs = map[string]map[string]string{
	"dev": {
		"APP_NAME": "myapp",
		"DEBUG":    "true",
		"PORT":     "8080",
	},
	"prod": {
		"APP_NAME": "myapp",
		"DEBUG":    "false",
		"PORT":     "443",
		"SENTRY":   "https://sentry.io/x",
	},
}

func TestBuild_AllKeysPresent(t *testing.T) {
	opts := DefaultOptions()
	rows := Build(envs, opts)
	if len(rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(rows))
	}
}

func TestBuild_ExcludeMissing(t *testing.T) {
	opts := Options{IncludeMissing: false}
	rows := Build(envs, opts)
	for _, r := range rows {
		if len(r.Missing) > 0 {
			t.Errorf("row %q should be excluded (missing in %v)", r.Key, r.Missing)
		}
	}
}

func TestRow_AllEqual_True(t *testing.T) {
	row := Row{
		Key:    "APP_NAME",
		Values: map[string]string{"dev": "myapp", "prod": "myapp"},
	}
	if !row.AllEqual() {
		t.Error("expected AllEqual to be true")
	}
}

func TestRow_AllEqual_False_Mismatch(t *testing.T) {
	row := Row{
		Key:    "DEBUG",
		Values: map[string]string{"dev": "true", "prod": "false"},
	}
	if row.AllEqual() {
		t.Error("expected AllEqual to be false due to mismatch")
	}
}

func TestRow_AllEqual_False_Missing(t *testing.T) {
	row := Row{
		Key:     "SENTRY",
		Values:  map[string]string{"dev": "", "prod": "https://sentry.io/x"},
		Missing: []string{"dev"},
	}
	if row.AllEqual() {
		t.Error("expected AllEqual to be false due to missing key")
	}
}

func TestFormatText_ContainsHeaders(t *testing.T) {
	opts := DefaultOptions()
	rows := Build(envs, opts)
	out := FormatText(rows, []string{"dev", "prod"})
	if !strings.Contains(out, "KEY") {
		t.Error("expected header to contain KEY")
	}
	if !strings.Contains(out, "dev") {
		t.Error("expected header to contain dev")
	}
}

func TestFormatText_EmptyRows(t *testing.T) {
	out := FormatText(nil, []string{"dev", "prod"})
	if !strings.Contains(out, "empty") {
		t.Errorf("expected empty matrix message, got: %s", out)
	}
}

func TestFormatText_StatusDIFF(t *testing.T) {
	opts := DefaultOptions()
	rows := Build(envs, opts)
	out := FormatText(rows, []string{"dev", "prod"})
	if !strings.Contains(out, "DIFF") {
		t.Error("expected DIFF status for mismatched keys")
	}
}
