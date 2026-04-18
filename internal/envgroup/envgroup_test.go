package envgroup

import (
	"testing"
)

func TestApply_GroupsByPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "envdiff",
	}
	groups := Apply(env, DefaultOptions())
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "APP" || groups[1].Name != "DB" {
		t.Errorf("unexpected group names: %v, %v", groups[0].Name, groups[1].Name)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	opts := DefaultOptions()
	opts.StripPrefix = true
	groups := Apply(env, opts)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	g := groups[0]
	if _, ok := g.Keys["HOST"]; !ok {
		t.Error("expected key HOST after strip")
	}
	if _, ok := g.Keys["PORT"]; !ok {
		t.Error("expected key PORT after strip")
	}
}

func TestApply_NoSeparator_EmptyGroup(t *testing.T) {
	env := map[string]string{"PLAINKEY": "value"}
	groups := Apply(env, DefaultOptions())
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "" {
		t.Errorf("expected empty group name, got %q", groups[0].Name)
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	groups := Apply(map[string]string{}, DefaultOptions())
	if len(groups) != 0 {
		t.Errorf("expected no groups, got %d", len(groups))
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	env := map[string]string{"DB.HOST": "localhost", "DB.PORT": "5432"}
	opts := Options{Separator: ".", StripPrefix: false}
	groups := Apply(env, opts)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "DB" {
		t.Errorf("expected group DB, got %q", groups[0].Name)
	}
}
