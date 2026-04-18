// Package envpatch applies a set of key-value patch operations to an env map.
// Supported operations: set, unset, rename.
package envpatch

import "fmt"

// OpType represents the type of patch operation.
type OpType string

const (
	OpSet    OpType = "set"
	OpUnset  OpType = "unset"
	OpRename OpType = "rename"
)

// Op describes a single patch operation.
type Op struct {
	Type  OpType
	Key   string
	Value string // used by OpSet
	To    string // used by OpRename
}

// Result summarises what changed after Apply.
type Result struct {
	Applied []string
	Skipped []string
}

// Apply executes patch ops against a copy of env and returns the patched map
// along with a Result describing what happened.
func Apply(env map[string]string, ops []Op) (map[string]string, Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var res Result
	for _, op := range ops {
		switch op.Type {
		case OpSet:
			out[op.Key] = op.Value
			res.Applied = append(res.Applied, fmt.Sprintf("set %s", op.Key))
		case OpUnset:
			if _, ok := out[op.Key]; !ok {
				res.Skipped = append(res.Skipped, fmt.Sprintf("unset %s (not found)", op.Key))
				continue
			}
			delete(out, op.Key)
			res.Applied = append(res.Applied, fmt.Sprintf("unset %s", op.Key))
		case OpRename:
			if op.To == "" {
				return nil, Result{}, fmt.Errorf("rename op for %q has empty 'to' field", op.Key)
			}
			val, ok := out[op.Key]
			if !ok {
				res.Skipped = append(res.Skipped, fmt.Sprintf("rename %s (not found)", op.Key))
				continue
			}
			delete(out, op.Key)
			out[op.To] = val
			res.Applied = append(res.Applied, fmt.Sprintf("rename %s -> %s", op.Key, op.To))
		default:
			return nil, Result{}, fmt.Errorf("unknown op type: %q", op.Type)
		}
	}
	return out, res, nil
}
