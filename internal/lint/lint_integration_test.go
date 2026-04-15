package lint_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/lint"
	"github.com/user/envdiff/internal/loader"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestLint_FullPipeline_NoIssues(t *testing.T) {
	path := writeTempEnvFile(t, "APP_ENV=production\nPORT=443\n")

	env, err := loader.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	issues := lint.Run(env, lint.DefaultOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got: %s", lint.FormatIssues(issues))
	}
}

func TestLint_FullPipeline_DetectsEmptyValue(t *testing.T) {
	path := writeTempEnvFile(t, "SECRET=\nHOST=localhost\n")

	env, err := loader.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	issues := lint.Run(env, lint.DefaultOptions())
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for empty SECRET")
	}

	found := false
	for _, iss := range issues {
		if iss.Key == "SECRET" && strings.Contains(iss.Message, "empty") {
			found = true
		}
	}
	if !found {
		t.Errorf("did not find expected issue for SECRET; issues: %s", lint.FormatIssues(issues))
	}
}

func TestLint_FullPipeline_LongValue(t *testing.T) {
	long := strings.Repeat("a", 600)
	path := writeTempEnvFile(t, "TOKEN="+long+"\n")

	env, err := loader.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	issues := lint.Run(env, lint.DefaultOptions())
	if len(issues) == 0 {
		t.Fatal("expected issue for oversized TOKEN value")
	}
	if !strings.Contains(lint.FormatIssues(issues), "exceeds maximum") {
		t.Errorf("expected 'exceeds maximum' in output")
	}
}
