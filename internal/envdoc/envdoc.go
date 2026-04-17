// Package envdoc extracts inline comments from .env files as documentation.
package envdoc

import (
	"bufio"
	"os"
	"strings"
)

// Entry holds a key, its value, and any associated comment.
type Entry struct {
	Key     string
	Value   string
	Comment string
}

// Parse reads a .env file and returns entries with inline or preceding comments.
func Parse(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	var pendingComment string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			pendingComment = ""
			continue
		}

		if strings.HasPrefix(line, "#") {
			pendingComment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			pendingComment = ""
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Inline comment after value
		inline := ""
		if idx := strings.Index(val, " #"); idx != -1 {
			inline = strings.TrimSpace(val[idx+2:])
			val = strings.TrimSpace(val[:idx])
		}

		comment := pendingComment
		if inline != "" {
			comment = inline
		}

		entries = append(entries, Entry{Key: key, Value: val, Comment: comment})
		pendingComment = ""
	}

	return entries, scanner.Err()
}

// ToMap converts entries to a map of key -> comment (keys without comments are included with empty string).
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Comment
	}
	return m
}
