package envpivot_test

import (
	"testing"

	"github.com/user/envdiff/internal/envdiff/envpivot"
)

func TestPivot_BasicRows(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "3000"},
		"prod": {"HOST": "example.com", "PORT": "443"},
	}
	rows := envpivot.Pivot(envs, envpivot.DefaultOptions())
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Key != "HOST" {
		t.Errorf("expected first key HOST, got %s", rows[0].Key)
	}
	if rows[0].Values["dev"] != "localhost" {
		t.Errorf("unexpected dev value: %s", rows[0].Values["dev"])
	}
}

func TestPivot_MissingKeyIncluded(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "DEBUG": "true"},
		"prod": {"HOST": "example.com"},
	}
	opts := envpivot.Options{IncludeAbsent: true}
	rows := envpivot.Pivot(envs, opts)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	// DEBUG row should have empty value for prod
	var debugRow *envpivot.Row
	for i := range rows {
		if rows[i].Key == "DEBUG" {
			debugRow = &rows[i]
		}
	}
	if debugRow == nil {
		t.Fatal("expected DEBUG row")
	}
	if debugRow.Values["prod"] != "" {
		t.Errorf("expected empty prod value, got %q", debugRow.Values["prod"])
	}
}

func TestPivot_MissingKeyExcluded(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "DEBUG": "true"},
		"prod": {"HOST": "example.com"},
	}
	opts := envpivot.Options{IncludeAbsent: false}
	rows := envpivot.Pivot(envs, opts)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row (only HOST), got %d", len(rows))
	}
	if rows[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", rows[0].Key)
	}
}

func TestRow_AllEqual_True(t *testing.T) {
	r := envpivot.Row{
		Key:    "PORT",
		Values: map[string]string{"dev": "8080", "prod": "8080"},
	}
	if !r.AllEqual([]string{"dev", "prod"}) {
		t.Error("expected all equal")
	}
}

func TestRow_AllEqual_False(t *testing.T) {
	r := envpivot.Row{
		Key:    "PORT",
		Values: map[string]string{"dev": "3000", "prod": "443"},
	}
	if r.AllEqual([]string{"dev", "prod"}) {
		t.Error("expected not equal")
	}
}

func TestPivot_EmptyEnvs(t *testing.T) {
	rows := envpivot.Pivot(map[string]map[string]string{}, envpivot.DefaultOptions())
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}
