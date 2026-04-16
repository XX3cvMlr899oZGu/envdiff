package envtype_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/envtype"
	"github.com/yourorg/envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestInferAll_FullPipeline(t *testing.T) {
	path := writeTempEnvFile(t, `
PORT=3000
DEBUG=false
API_URL=https://api.example.com
ADMIN_EMAIL=admin@example.com
DATA_DIR=/var/data
APP_NAME=myapp
RATE=0.75
`)

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	types := envtype.InferAll(env)

	cases := map[string]envtype.Type{
		"PORT":        envtype.TypeInt,
		"DEBUG":       envtype.TypeBool,
		"API_URL":     envtype.TypeURL,
		"ADMIN_EMAIL": envtype.TypeEmail,
		"DATA_DIR":    envtype.TypePath,
		"APP_NAME":    envtype.TypeString,
		"RATE":        envtype.TypeFloat,
	}

	for key, want := range cases {
		if got := types[key]; got != want {
			t.Errorf("key %q: got type %q, want %q", key, got, want)
		}
	}
}
