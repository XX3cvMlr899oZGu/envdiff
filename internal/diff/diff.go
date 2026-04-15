package diff

// Result holds the comparison result between two env files.
type Result struct {
	// MissingInSecond contains keys present in first but absent in second.
	MissingInSecond []string
	// MissingInFirst contains keys present in second but absent in first.
	MissingInFirst []string
	// Mismatched contains keys present in both files but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey represents a key whose value differs between two env files.
type MismatchedKey struct {
	Key    string
	First  string
	Second string
}

// Compare compares two parsed env maps and returns a Result describing
// missing and mismatched keys.
func Compare(first, second map[string]string) Result {
	result := Result{}

	for key, val := range first {
		if secondVal, ok := second[key]; !ok {
			result.MissingInSecond = append(result.MissingInSecond, key)
		} else if val != secondVal {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:    key,
				First:  val,
				Second: secondVal,
			})
		}
	}

	for key := range second {
		if _, ok := first[key]; !ok {
			result.MissingInFirst = append(result.MissingInFirst, key)
		}
	}

	return result
}

// HasDifferences returns true if the Result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInFirst) > 0 ||
		len(r.MissingInSecond) > 0 ||
		len(r.Mismatched) > 0
}
