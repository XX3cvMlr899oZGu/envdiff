// Package envmatrix builds a comparison matrix across multiple env maps,
// showing each key's value per environment and whether they all agree.
package envmatrix

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultOptions returns a default Options value.
func DefaultOptions() Options {
	return Options{
		IncludeMissing: true,
	}
}

// Options controls matrix generation behaviour.
type Options struct {
	// IncludeMissing includes rows where a key is absent in some environments.
	IncludeMissing bool
}

// Row represents a single key across all environments.
type Row struct {
	Key    string
	Values map[string]string // env name → value (empty string means absent)
	Missing []string          // env names where the key is absent
}

// AllEqual returns true when every environment has the same non-empty value.
func (r Row) AllEqual() bool {
	if len(r.Missing) > 0 {
		return false
	}
	var first string
	for _, v := range r.Values {
		first = v
		break
	}
	for _, v := range r.Values {
		if v != first {
			return false
		}
	}
	return true
}

// Build constructs a matrix from a map of environment name → key/value pairs.
func Build(envs map[string]map[string]string, opts Options) []Row {
	keys := collectKeys(envs)
	envNames := sortedEnvNames(envs)

	var rows []Row
	for _, key := range keys {
		row := Row{
			Key:    key,
			Values: make(map[string]string, len(envNames)),
		}
		for _, name := range envNames {
			val, ok := envs[name][key]
			if ok {
				row.Values[name] = val
			} else {
				row.Missing = append(row.Missing, name)
				row.Values[name] = ""
			}
		}
		if !opts.IncludeMissing && len(row.Missing) > 0 {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

// FormatText renders the matrix as a plain-text table.
func FormatText(rows []Row, envNames []string) string {
	if len(rows) == 0 {
		return "(empty matrix)\n"
	}
	var sb strings.Builder
	header := fmt.Sprintf("%-30s", "KEY")
	for _, n := range envNames {
		header += fmt.Sprintf(" %-20s", n)
	}
	header += " STATUS"
	sb.WriteString(header + "\n")
	sb.WriteString(strings.Repeat("-", len(header)) + "\n")
	for _, row := range rows {
		line := fmt.Sprintf("%-30s", row.Key)
		for _, n := range envNames {
			line += fmt.Sprintf(" %-20s", row.Values[n])
		}
		status := "OK"
		if !row.AllEqual() {
			status = "DIFF"
		}
		sb.WriteString(line + " " + status + "\n")
	}
	return sb.String()
}

func collectKeys(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	for _, m := range envs {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedEnvNames(envs map[string]map[string]string) []string {
	names := make([]string, 0, len(envs))
	for n := range envs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
