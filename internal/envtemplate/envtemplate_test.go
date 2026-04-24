package envtemplate

import (
	"testing"
)

func TestApply_NoTemplates(t *testing.T) {
	tmplMap := map[string]string{"HOST": "localhost", "PORT": "5432"}
	env := map[string]string{}
	out, err := Apply(tmplMap, env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "5432" {
		t.Errorf("expected passthrough values, got %v", out)
	}
}

func TestApply_SimpleSubstitution(t *testing.T) {
	tmplMap := map[string]string{"DSN": "postgres://{{.DB_USER}}:{{.DB_PASS}}@localhost/mydb"}
	env := map[string]string{"DB_USER": "admin", "DB_PASS": "secret"}
	out, err := Apply(tmplMap, env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://admin:secret@localhost/mydb"
	if out["DSN"] != want {
		t.Errorf("got %q, want %q", out["DSN"], want)
	}
}

func TestApply_MissingKey_Error(t *testing.T) {
	tmplMap := map[string]string{"URL": "http://{{.MISSING_HOST}}/path"}
	env := map[string]string{}
	_, err := Apply(tmplMap, env, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestApply_MissingKey_Zero(t *testing.T) {
	opts := Options{MissingKey: "zero"}
	tmplMap := map[string]string{"URL": "http://{{.MISSING_HOST}}/path"}
	env := map[string]string{}
	out, err := Apply(tmplMap, env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http:///path" {
		t.Errorf("unexpected output: %q", out["URL"])
	}
}

func TestRenderString_Basic(t *testing.T) {
	env := map[string]string{"APP_ENV": "production"}
	result, err := RenderString("env={{.APP_ENV}}", env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "env=production" {
		t.Errorf("got %q, want %q", result, "env=production")
	}
}

func TestRenderString_InvalidTemplate(t *testing.T) {
	env := map[string]string{}
	_, err := RenderString("{{.Unclosed", env, DefaultOptions())
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestApply_MultipleKeys(t *testing.T) {
	tmplMap := map[string]string{
		"A": "hello-{{.NAME}}",
		"B": "static",
		"C": "{{.X}}-{{.Y}}",
	}
	env := map[string]string{"NAME": "world", "X": "foo", "Y": "bar"}
	out, err := Apply(tmplMap, env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "hello-world" {
		t.Errorf("A: got %q", out["A"])
	}
	if out["B"] != "static" {
		t.Errorf("B: got %q", out["B"])
	}
	if out["C"] != "foo-bar" {
		t.Errorf("C: got %q", out["C"])
	}
}
