// Package envtag provides functionality to tag env keys with arbitrary
// labels and filter or group entries by those tags.
package envtag

import (
	"fmt"
	"strings"
)

// Tag represents a label attached to one or more env keys.
type Tag struct {
	Name string
	Keys []string
}

// Options controls tagging behaviour.
type Options struct {
	// Tags is the list of tag definitions to apply.
	Tags []Tag
}

// DefaultOptions returns a zero-value Options.
func DefaultOptions() Options {
	return Options{}
}

// Result maps each tag name to the subset of env entries whose keys carry that tag.
type Result map[string]map[string]string

// Apply scans env and returns a Result grouping keys by their assigned tags.
// A key may appear in multiple tag groups. Keys with no matching tag are
// placed under the reserved "_untagged" group.
func Apply(env map[string]string, opts Options) (Result, error) {
	if env == nil {
		return Result{}, nil
	}

	// Build a lookup: key -> list of tag names
	keyTags := make(map[string][]string)
	for _, tag := range opts.Tags {
		if strings.TrimSpace(tag.Name) == "" {
			return nil, fmt.Errorf("envtag: tag name must not be empty")
		}
		for _, k := range tag.Keys {
			keyTags[k] = append(keyTags[k], tag.Name)
		}
	}

	res := Result{}
	for k, v := range env {
		tags, ok := keyTags[k]
		if !ok {
			if res["_untagged"] == nil {
				res["_untagged"] = map[string]string{}
			}
			res["_untagged"][k] = v
			continue
		}
		for _, t := range tags {
			if res[t] == nil {
				res[t] = map[string]string{}
			}
			res[t][k] = v
		}
	}
	return res, nil
}

// KeysForTag returns all keys in env that are associated with tagName.
func KeysForTag(result Result, tagName string) []string {
	group, ok := result[tagName]
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(group))
	for k := range group {
		keys = append(keys, k)
	}
	return keys
}
