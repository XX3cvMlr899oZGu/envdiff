package diff

import "sort"

// Status represents the comparison result for a single key.
type Status string

const (
	StatusEqual    Status = "equal"
	StatusMissing  Status = "missing"
	StatusExtra    Status = "extra"
	StatusMismatch Status = "mismatch"
)

// Result holds the comparison outcome for a single environment key.
type Result struct {
	Key    string
	Status Status
	ValueA string
	ValueB string
}

// Compare compares two env maps and returns a sorted slice of Results.
func Compare(a, b map[string]string) []Result {
	keys := make(map[string]struct{})
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	results := make([]Result, 0, len(keys))
	for k := range keys {
		va, inA := a[k]
		vb, inB := b[k]

		var status Status
		switch {
		case inA && !inB:
			status = StatusMissing
		case !inA && inB:
			status = StatusExtra
		case va == vb:
			status = StatusEqual
		default:
			status = StatusMismatch
		}

		results = append(results, Result{Key: k, Status: status, ValueA: va, ValueB: vb})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}
