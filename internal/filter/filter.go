package filter

import (
	"regexp"
	"strings"
)

// Options holds filtering configuration.
type Options struct {
	Prefix    string
	KeyRegex  string
	Exclude   []string
}

// ApplyToMap returns a new map containing only the entries that pass
// the filter defined by opts.
func ApplyToMap(env map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.KeyRegex != "" {
		var err error
		re, err = regexp.Compile(opts.KeyRegex)
		if err != nil {
			return nil, err
		}
	}

	excludeSet := make(map[string]struct{}, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[k] = struct{}{}
	}

	out := make(map[string]string)
	for k, v := range env {
		if _, excluded := excludeSet[k]; excluded {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		out[k] = v
	}
	return out, nil
}
