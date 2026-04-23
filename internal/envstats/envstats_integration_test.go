package envstats_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/envstats"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/filter"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envstats-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestEnvStats_FullPipeline_Basic(t *testing.T) {
	path := writeTempEnvFile(t, `
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=supersecretpassword
APP_ENV=production
APP_DEBUG=
`)

	env, err := loader.LoadFile(path, filter.Options{})
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	s := envstats.Compute(env)

	if s.TotalKeys != 5 {
		t.Errorf("expected 5 keys, got %d", s.TotalKeys)
	}
	if s.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", s.EmptyValues)
	}
	if s.UppercaseKeys != 5 {
		t.Errorf("expected all 5 keys to be uppercase, got %d", s.UppercaseKeys)
	}
	if s.MaxValueLength < len("supersecretpassword") {
		t.Errorf("expected max value length >= %d, got %d", len("supersecretpassword"), s.MaxValueLength)
	}
}

func TestEnvStats_FullPipeline_TopLongest(t *testing.T) {
	path := writeTempEnvFile(t, `
SHORT=hi
MEDIUM=mediumvalue
LONG=this_is_a_very_long_environment_variable_value
`)

	env, err := loader.LoadFile(path, filter.Options{})
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	top := envstats.TopLongestValues(env, 1)
	if len(top) != 1 {
		t.Fatalf("expected 1 result, got %d", len(top))
	}
	if top[0] != "LONG" {
		t.Errorf("expected LONG, got %s", top[0])
	}
}
