package diff

// Status describes the relationship of a key between two env files.
type Status int

const (
	StatusEqual      Status = iota
	StatusMissingInA        // key exists in B but not A
	StatusMissingInB        // key exists in A but not B
	StatusMismatch          // key exists in both but values differ
)

// Result holds the comparison outcome for a single key.
type Result struct {
	Key    string
	Status Status
	ValueA string
	ValueB string
}

// Compare compares two maps of env key/value pairs and returns a slice of Results.
func Compare(a, b map[string]string) []Result {
	var results []Result

	for k, va := range a {
		if vb, ok := b[k]; !ok {
			results = append(results, Result{Key: k, Status: StatusMissingInB, ValueA: va})
		} else if va != vb {
			results = append(results, Result{Key: k, Status: StatusMismatch, ValueA: va, ValueB: vb})
		} else {
			results = append(results, Result{Key: k, Status: StatusEqual, ValueA: va, ValueB: vb})
		}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			results = append(results, Result{Key: k, Status: StatusMissingInA, ValueB: vb})
		}
	}

	return results
}
