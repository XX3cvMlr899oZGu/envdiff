package envtype_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/envtype"
)

func TestInfer_Bool(t *testing.T) {
	for _, v := range []string{"true", "false", "TRUE", "1", "0"} {
		if got := envtype.Infer(v); got != envtype.TypeBool {
			t.Errorf("Infer(%q) = %q, want bool", v, got)
		}
	}
}

func TestInfer_Int(t *testing.T) {
	for _, v := range []string{"42", "-7", "1000"} {
		if got := envtype.Infer(v); got != envtype.TypeInt {
			t.Errorf("Infer(%q) = %q, want int", v, got)
		}
	}
}

func TestInfer_Float(t *testing.T) {
	for _, v := range []string{"3.14", "-0.5", "2.0"} {
		if got := envtype.Infer(v); got != envtype.TypeFloat {
			t.Errorf("Infer(%q) = %q, want float", v, got)
		}
	}
}

func TestInfer_URL(t *testing.T) {
	for _, v := range []string{"http://example.com", "https://api.example.com/v1"} {
		if got := envtype.Infer(v); got != envtype.TypeURL {
			t.Errorf("Infer(%q) = %q, want url", v, got)
		}
	}
}

func TestInfer_Email(t *testing.T) {
	if got := envtype.Infer("user@example.com"); got != envtype.TypeEmail {
		t.Errorf("Infer(email) = %q, want email", got)
	}
}

func TestInfer_Path(t *testing.T) {
	for _, v := range []string{"/etc/app", "~/config", "./local"} {
		if got := envtype.Infer(v); got != envtype.TypePath {
			t.Errorf("Infer(%q) = %q, want path", v, got)
		}
	}
}

func TestInfer_String(t *testing.T) {
	if got := envtype.Infer("hello world"); got != envtype.TypeString {
		t.Errorf("Infer(string) = %q, want string", got)
	}
}

func TestInfer_Empty(t *testing.T) {
	if got := envtype.Infer(""); got != envtype.TypeString {
		t.Errorf("Infer(\"\") = %q, want string", got)
	}
}

func TestInferAll(t *testing.T) {
	env := map[string]string{
		"PORT":     "8080",
		"DEBUG":    "true",
		"BASE_URL": "https://example.com",
		"LABEL":    "production",
	}
	result := envtype.InferAll(env)
	if result["PORT"] != envtype.TypeInt {
		t.Errorf("PORT: got %q, want int", result["PORT"])
	}
	if result["DEBUG"] != envtype.TypeBool {
		t.Errorf("DEBUG: got %q, want bool", result["DEBUG"])
	}
	if result["BASE_URL"] != envtype.TypeURL {
		t.Errorf("BASE_URL: got %q, want url", result["BASE_URL"])
	}
	if result["LABEL"] != envtype.TypeString {
		t.Errorf("LABEL: got %q, want string", result["LABEL"])
	}
}
