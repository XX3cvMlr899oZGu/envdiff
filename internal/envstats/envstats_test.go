package envstats_test

import (
	"testing"

	"github.com/user/envdiff/internal/envstats"
)

func TestCompute_EmptyMap(t *testing.T) {
	s := envstats.Compute(map[string]string{})
	if s.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", s.TotalKeys)
	}
}

func TestCompute_BasicCounts(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"SECRET":  "",
	}
	s := envstats.Compute(env)
	if s.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", s.TotalKeys)
	}
	if s.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", s.EmptyValues)
	}
	if s.UppercaseKeys != 3 {
		t.Errorf("expected 3 uppercase keys, got %d", s.UppercaseKeys)
	}
}

func TestCompute_MixedCaseKeys(t *testing.T) {
	env := map[string]string{
		"MyKey":   "val1",
		"another": "val2",
		"UPPER":   "val3",
	}
	s := envstats.Compute(env)
	if s.MixedCaseKeys != 1 {
		t.Errorf("expected 1 mixed-case key, got %d", s.MixedCaseKeys)
	}
	if s.LowercaseKeys != 1 {
		t.Errorf("expected 1 lowercase key, got %d", s.LowercaseKeys)
	}
	if s.UppercaseKeys != 1 {
		t.Errorf("expected 1 uppercase key, got %d", s.UppercaseKeys)
	}
}

func TestCompute_ValueLengths(t *testing.T) {
	env := map[string]string{
		"A": "hi",
		"B": "hello world",
		"C": "x",
	}
	s := envstats.Compute(env)
	if s.MaxValueLength != 11 {
		t.Errorf("expected max 11, got %d", s.MaxValueLength)
	}
	if s.MinValueLength != 1 {
		t.Errorf("expected min 1, got %d", s.MinValueLength)
	}
	expectedAvg := float64(2+11+1) / 3.0
	if s.AvgValueLength != expectedAvg {
		t.Errorf("expected avg %.4f, got %.4f", expectedAvg, s.AvgValueLength)
	}
}

func TestCompute_UniqueValues(t *testing.T) {
	env := map[string]string{
		"A": "same",
		"B": "same",
		"C": "different",
	}
	s := envstats.Compute(env)
	if s.UniqueValues != 2 {
		t.Errorf("expected 2 unique values, got %d", s.UniqueValues)
	}
}

func TestTopLongestValues_Order(t *testing.T) {
	env := map[string]string{
		"SHORT": "hi",
		"LONG":  "this is a long value",
		"MED":   "medium",
	}
	top := envstats.TopLongestValues(env, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0] != "LONG" {
		t.Errorf("expected LONG first, got %s", top[0])
	}
}

func TestTopLongestValues_NLargerThanMap(t *testing.T) {
	env := map[string]string{"A": "one", "B": "two"}
	top := envstats.TopLongestValues(env, 10)
	if len(top) != 2 {
		t.Errorf("expected 2, got %d", len(top))
	}
}
