package envpatch

import (
	"testing"
)

func base() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"SECRET":   "abc123",
	}
}

func TestApply_SetNewKey(t *testing.T) {
	out, res, err := Apply(base(), []Op{{Type: OpSet, Key: "NEW_KEY", Value: "hello"}})
	if err != nil { t.Fatal(err) }
	if out["NEW_KEY"] != "hello" { t.Errorf("expected hello, got %s", out["NEW_KEY"]) }
	if len(res.Applied) != 1 { t.Errorf("expected 1 applied, got %d", len(res.Applied)) }
}

func TestApply_SetOverwritesExisting(t *testing.T) {
	out, _, err := Apply(base(), []Op{{Type: OpSet, Key: "APP_ENV", Value: "staging"}})
	if err != nil { t.Fatal(err) }
	if out["APP_ENV"] != "staging" { t.Errorf("expected staging, got %s", out["APP_ENV"]) }
}

func TestApply_UnsetExistingKey(t *testing.T) {
	out, res, err := Apply(base(), []Op{{Type: OpUnset, Key: "SECRET"}})
	if err != nil { t.Fatal(err) }
	if _, ok := out["SECRET"]; ok { t.Error("expected SECRET to be removed") }
	if len(res.Applied) != 1 { t.Errorf("expected 1 applied") }
}

func TestApply_UnsetMissingKey_Skipped(t *testing.T) {
	_, res, err := Apply(base(), []Op{{Type: OpUnset, Key: "MISSING"}})
	if err != nil { t.Fatal(err) }
	if len(res.Skipped) != 1 { t.Errorf("expected 1 skipped, got %d", len(res.Skipped)) }
}

func TestApply_RenameKey(t *testing.T) {
	out, res, err := Apply(base(), []Op{{Type: OpRename, Key: "DB_HOST", To: "DATABASE_HOST"}})
	if err != nil { t.Fatal(err) }
	if _, ok := out["DB_HOST"]; ok { t.Error("old key should be gone") }
	if out["DATABASE_HOST"] != "localhost" { t.Errorf("expected localhost") }
	if len(res.Applied) != 1 { t.Errorf("expected 1 applied") }
}

func TestApply_RenameMissingKey_Skipped(t *testing.T) {
	_, res, err := Apply(base(), []Op{{Type: OpRename, Key: "GHOST", To: "SPIRIT"}})
	if err != nil { t.Fatal(err) }
	if len(res.Skipped) != 1 { t.Errorf("expected 1 skipped") }
}

func TestApply_RenameEmptyTo_Error(t *testing.T) {
	_, _, err := Apply(base(), []Op{{Type: OpRename, Key: "DB_HOST", To: ""}})
	if err == nil { t.Error("expected error for empty 'to'") }
}

func TestApply_UnknownOp_Error(t *testing.T) {
	_, _, err := Apply(base(), []Op{{Type: "copy", Key: "X"}})
	if err == nil { t.Error("expected error for unknown op") }
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := base()
	Apply(env, []Op{{Type: OpUnset, Key: "SECRET"}})
	if _, ok := env["SECRET"]; !ok { t.Error("original map should not be mutated") }
}
