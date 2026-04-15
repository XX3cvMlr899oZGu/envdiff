package main

import (
	"flag"
	"strings"
)

// patternList is a flag.Value that accumulates multiple --redact-pattern values.
type patternList []string

func (p *patternList) String() string {
	if p == nil || len(*p) == 0 {
		return ""
	}
	return strings.Join(*p, ",")
}

func (p *patternList) Set(val string) error {
	*p = append(*p, val)
	return nil
}

var (
	redactFlag      bool
	redactPatterns  patternList
	redactMask      string
)

func registerRedactFlags(fs *flag.FlagSet) {
	fs.BoolVar(&redactFlag, "redact", false, "mask sensitive values before output")
	fs.Var(&redactPatterns, "redact-pattern", "additional regex pattern for sensitive key detection (repeatable)")
	fs.StringVar(&redactMask, "redact-mask", "***REDACTED***", "placeholder string for redacted values")
}
