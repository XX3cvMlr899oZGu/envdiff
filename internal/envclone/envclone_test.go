package envclone

import (
	"strings"
	"testing"
)

func base() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db",
		"DB_PORT":  "5432",
	}
}

func TestApply_NoOptions(t *testing.T) {
	src := base()
	out := Apply(src, DefaultOptions())
	if len(out) != len(src) {
		t.Fatalf("expected %d keys, got %d", len(src), len(out))
	}
	for k, v := range src {
		if out[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	out := Apply(base(), Options{KeyPrefix: "APP_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should be excluded")
	}
}

func TestApply_StripPrefix(t *testing.T) {
	out := Apply(base(), Options{KeyPrefix: "APP_", StripPrefix: true})
	if _, ok := out["HOST"]; !ok {
		t.Error("expected key HOST after stripping prefix")
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("APP_HOST should not appear after stripping")
	}
}

func TestApply_KeyTransform(t *testing.T) {
	out := Apply(base(), Options{KeyTransform: strings.ToLower})
	if _, ok := out["app_host"]; !ok {
		t.Error("expected lowercase key app_host")
	}
}

func TestApply_DoesNotMutateSource(t *testing.T) {
	src := base()
	out := Apply(src, Options{KeyTransform: strings.ToLower})
	_ = out
	if _, ok := src["APP_HOST"]; !ok {
		t.Error("source map was mutated")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	out := Apply(map[string]string{}, Options{KeyPrefix: "APP_"})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}
