package envcount_test

import (
	"testing"

	"github.com/user/envdiff/internal/envcount"
)

func TestApply_BasicCounts(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "myapp",
		"EMPTY_KEY": "",
	}

	s := envcount.Apply(env, envcount.DefaultOptions())

	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
	if s.NonEmpty != 3 {
		t.Errorf("expected NonEmpty=3, got %d", s.NonEmpty)
	}
}

func TestApply_PrefixCounts(t *testing.T) {
	env := map[string]string{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
		"DB_NAME":   "mydb",
		"APP_NAME":  "myapp",
		"NOPREFIXKEY": "val",
	}

	s := envcount.Apply(env, envcount.DefaultOptions())

	if s.Prefixes["DB"] != 3 {
		t.Errorf("expected DB prefix count=3, got %d", s.Prefixes["DB"])
	}
	if s.Prefixes["APP"] != 1 {
		t.Errorf("expected APP prefix count=1, got %d", s.Prefixes["APP"])
	}
}

func TestApply_EmptyMap(t *testing.T) {
	s := envcount.Apply(map[string]string{}, envcount.DefaultOptions())

	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
	if len(s.Prefixes) != 0 {
		t.Errorf("expected no prefixes, got %v", s.Prefixes)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	env := map[string]string{
		"DB.HOST": "localhost",
		"DB.PORT": "5432",
		"APP.NAME": "myapp",
	}

	opts := envcount.Options{Separator: "."}
	s := envcount.Apply(env, opts)

	if s.Prefixes["DB"] != 2 {
		t.Errorf("expected DB prefix count=2, got %d", s.Prefixes["DB"])
	}
}

func TestTopPrefixes_Ordered(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "h",
		"DB_PORT":  "p",
		"DB_NAME":  "n",
		"APP_NAME": "a",
		"APP_ENV":  "e",
		"SVC_URL":  "u",
	}

	s := envcount.Apply(env, envcount.DefaultOptions())
	top := envcount.TopPrefixes(s, 2)

	if len(top) != 2 {
		t.Fatalf("expected 2 top prefixes, got %d", len(top))
	}
	if top[0] != "DB" {
		t.Errorf("expected first prefix=DB, got %s", top[0])
	}
	if top[1] != "APP" {
		t.Errorf("expected second prefix=APP, got %s", top[1])
	}
}

func TestTopPrefixes_AllWhenNZero(t *testing.T) {
	env := map[string]string{
		"A_1": "v",
		"B_1": "v",
		"C_1": "v",
	}

	s := envcount.Apply(env, envcount.DefaultOptions())
	top := envcount.TopPrefixes(s, 0)

	if len(top) != 3 {
		t.Errorf("expected 3 prefixes, got %d", len(top))
	}
}
