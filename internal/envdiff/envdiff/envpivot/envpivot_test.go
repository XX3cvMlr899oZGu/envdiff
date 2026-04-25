package envpivot

import (
	"testing"
)

func envs() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {"APP_NAME": "myapp", "DB_HOST": "localhost", "DEBUG": "true"},
		"prod": {"APP_NAME": "myapp", "DB_HOST": "prod.db"},
	}
}

func TestPivot_BasicRows(t *testing.T) {
	rows := Pivot(envs(), DefaultOptions())
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[0].Key != "APP_NAME" {
		t.Errorf("expected first key APP_NAME, got %s", rows[0].Key)
	}
}

func TestPivot_MissingKeyIncluded(t *testing.T) {
	rows := Pivot(envs(), Options{ExcludeMissing: false})
	var debugRow *Row
	for i := range rows {
		if rows[i].Key == "DEBUG" {
			debugRow = &rows[i]
			break
		}
	}
	if debugRow == nil {
		t.Fatal("expected DEBUG row to be present")
	}
	if _, ok := debugRow.Values["dev"]; !ok {
		t.Error("expected DEBUG to have dev value")
	}
	if _, ok := debugRow.Values["prod"]; ok {
		t.Error("expected DEBUG to be absent in prod")
	}
}

func TestPivot_MissingKeyExcluded(t *testing.T) {
	rows := Pivot(envs(), Options{ExcludeMissing: true})
	for _, r := range rows {
		if r.Key == "DEBUG" {
			t.Error("DEBUG should have been excluded")
		}
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
}

func TestRow_AllEqual_True(t *testing.T) {
	r := Row{Key: "APP_NAME", Values: map[string]string{"dev": "myapp", "prod": "myapp"}}
	if !r.AllEqual() {
		t.Error("expected AllEqual to be true")
	}
}

func TestRow_AllEqual_False(t *testing.T) {
	r := Row{Key: "DB_HOST", Values: map[string]string{"dev": "localhost", "prod": "prod.db"}}
	if r.AllEqual() {
		t.Error("expected AllEqual to be false")
	}
}

func TestPivot_EmptyEnvs(t *testing.T) {
	rows := Pivot(map[string]map[string]string{}, DefaultOptions())
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty input, got %d", len(rows))
	}
}

func TestPivot_SortedByKey(t *testing.T) {
	rows := Pivot(envs(), DefaultOptions())
	for i := 1; i < len(rows); i++ {
		if rows[i].Key < rows[i-1].Key {
			t.Errorf("rows not sorted: %s before %s", rows[i-1].Key, rows[i].Key)
		}
	}
}
