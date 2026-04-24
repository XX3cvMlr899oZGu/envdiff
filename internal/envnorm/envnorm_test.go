package envnorm_test

import (
	"testing"

	"github.com/user/envdiff/internal/envnorm"
)

func TestApply_UppercaseKeys(t *testing.T) {
	env := map[string]string{"app_name": "myapp", "db_host": "localhost"}
	opts := envnorm.DefaultOptions()
	out, err := envnorm.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", out["APP_NAME"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestApply_LowercaseKeys(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "DB_HOST": "localhost"}
	opts := envnorm.Options{Style: envnorm.StyleLower}
	out, err := envnorm.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["app_name"] != "myapp" {
		t.Errorf("expected app_name=myapp, got %q", out["app_name"])
	}
}

func TestApply_ReplaceHyphens(t *testing.T) {
	env := map[string]string{"app-name": "myapp"}
	opts := envnorm.Options{Style: envnorm.StyleUpper, ReplaceHyphens: true}
	out, err := envnorm.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %v", out)
	}
}

func TestApply_ReplaceDots(t *testing.T) {
	env := map[string]string{"db.host": "localhost"}
	opts := envnorm.Options{Style: envnorm.StyleUpper, ReplaceDots: true}
	out, err := envnorm.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
}

func TestApply_CollisionReturnsError(t *testing.T) {
	// "app-name" and "app_name" both normalize to "APP_NAME"
	env := map[string]string{"app-name": "a", "app_name": "b"}
	opts := envnorm.Options{Style: envnorm.StyleUpper, ReplaceHyphens: true}
	_, err := envnorm.Apply(env, opts)
	if err == nil {
		t.Error("expected collision error, got nil")
	}
}

func TestApply_StyleNone_NoChange(t *testing.T) {
	env := map[string]string{"MyKey": "val"}
	opts := envnorm.Options{Style: envnorm.StyleNone}
	out, err := envnorm.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MyKey"] != "val" {
		t.Errorf("expected MyKey unchanged, got %v", out)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	out, err := envnorm.Apply(map[string]string{}, envnorm.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
