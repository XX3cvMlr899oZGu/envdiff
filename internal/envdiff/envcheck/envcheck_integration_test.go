package envcheck_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/envdiff/envcheck"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/filter"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envcheck-*.env")
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

func TestEnvCheck_FullPipeline_AllPass(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=production\n")

	env, err := loader.LoadFile(path, filter.ApplyToMap)
	if err != nil {
		t.Fatalf("loader.LoadFile: %v", err)
	}

	rules := []envcheck.Rule{
		{Key: "DB_HOST", Required: true},
		{Key: "DB_PORT", Required: true},
		{Key: "APP_ENV", Required: true, NonEmpty: true},
	}

	violations := envcheck.Check(env, rules)
	if envcheck.HasViolations(violations) {
		t.Errorf("expected no violations, got: %v", envcheck.FormatViolations(violations))
	}
}

func TestEnvCheck_FullPipeline_RequiredMissing(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\n")

	env, err := loader.LoadFile(path, filter.ApplyToMap)
	if err != nil {
		t.Fatalf("loader.LoadFile: %v", err)
	}

	rules := []envcheck.Rule{
		{Key: "DB_HOST", Required: true},
		{Key: "DB_PORT", Required: true},
		{Key: "SECRET_KEY", Required: true, NonEmpty: true},
	}

	violations := envcheck.Check(env, rules)
	if !envcheck.HasViolations(violations) {
		t.Fatal("expected violations but got none")
	}

	formatted := envcheck.FormatViolations(violations)
	if len(formatted) == 0 {
		t.Error("expected non-empty formatted violations")
	}
}

func TestEnvCheck_FullPipeline_ForbiddenPresent(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\nDEBUG=true\n")

	env, err := loader.LoadFile(path, filter.ApplyToMap)
	if err != nil {
		t.Fatalf("loader.LoadFile: %v", err)
	}

	rules := []envcheck.Rule{
		{Key: "DEBUG", Forbidden: true},
	}

	violations := envcheck.Check(env, rules)
	if !envcheck.HasViolations(violations) {
		t.Fatal("expected violation for forbidden key DEBUG")
	}
}
