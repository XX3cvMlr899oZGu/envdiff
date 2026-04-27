package envlabel

import (
	"strings"
	"testing"
)

var baseEnv = map[string]string{
	"APP_HOST": "localhost",
	"APP_PORT": "8080",
}

func TestApply_NoLabels(t *testing.T) {
	opts := DefaultOptions()
	le, err := Apply(baseEnv, nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(le.Labels) != 0 {
		t.Errorf("expected 0 labels, got %d", len(le.Labels))
	}
	if le.Env["APP_HOST"] != "localhost" {
		t.Errorf("env not copied correctly")
	}
}

func TestApply_AddsLabels(t *testing.T) {
	opts := DefaultOptions()
	opts.Labels = map[string]string{"env": "production", "team": "platform"}
	le, err := Apply(baseEnv, nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if le.Labels["env"] != "production" {
		t.Errorf("expected label env=production")
	}
	if le.Labels["team"] != "platform" {
		t.Errorf("expected label team=platform")
	}
}

func TestApply_ConflictReturnsError(t *testing.T) {
	existing := map[string]string{"env": "staging"}
	opts := DefaultOptions()
	opts.Labels = map[string]string{"env": "production"}
	_, err := Apply(baseEnv, existing, opts)
	if err == nil {
		t.Fatal("expected error on label conflict")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestApply_OverwriteExisting(t *testing.T) {
	existing := map[string]string{"env": "staging"}
	opts := DefaultOptions()
	opts.Labels = map[string]string{"env": "production"}
	opts.OverwriteExisting = true
	le, err := Apply(baseEnv, existing, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if le.Labels["env"] != "production" {
		t.Errorf("expected overwritten label env=production")
	}
}

func TestApply_EmptyKeyReturnsError(t *testing.T) {
	opts := DefaultOptions()
	opts.Labels = map[string]string{"": "value"}
	_, err := Apply(baseEnv, nil, opts)
	if err == nil {
		t.Fatal("expected error for empty label key")
	}
}

func TestApply_NilEnvReturnsError(t *testing.T) {
	opts := DefaultOptions()
	_, err := Apply(nil, nil, opts)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestFormatText_NoLabels(t *testing.T) {
	le := LabeledEnv{Env: baseEnv, Labels: map[string]string{}}
	out := FormatText(le)
	if out != "(no labels)" {
		t.Errorf("expected '(no labels)', got %q", out)
	}
}

func TestFormatText_WithLabels(t *testing.T) {
	le := LabeledEnv{
		Env:    baseEnv,
		Labels: map[string]string{"env": "prod", "region": "us-east-1"},
	}
	out := FormatText(le)
	if !strings.Contains(out, "env=prod") {
		t.Errorf("expected env=prod in output")
	}
	if !strings.Contains(out, "region=us-east-1") {
		t.Errorf("expected region=us-east-1 in output")
	}
}
