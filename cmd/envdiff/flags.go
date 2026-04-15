package main

import (
	"flag"
	"strings"
)

// excludeList is a custom flag.Value that collects repeated --exclude flags.
type excludeList []string

func (e *excludeList) String() string {
	return strings.Join(*e, ",")
}

func (e *excludeList) Set(value string) error {
	*e = append(*e, value)
	return nil
}

// cliFlags holds all parsed command-line options.
type cliFlags struct {
	FileA    string
	FileB    string
	Format   string
	ShowAll  bool
	Prefix   string
	KeyRegex string
	Exclude  excludeList
}

func parseFlags() cliFlags {
	var f cliFlags

	flag.StringVar(&f.FileA, "a", "", "path to first .env file (required)")
	flag.StringVar(&f.FileB, "b", "", "path to second .env file (required)")
	flag.StringVar(&f.Format, "format", "text", "output format: text or json")
	flag.BoolVar(&f.ShowAll, "show-all", false, "include keys with equal values in output")
	flag.StringVar(&f.Prefix, "prefix", "", "only compare keys with this prefix")
	flag.StringVar(&f.KeyRegex, "key-regex", "", "only compare keys matching this regular expression")
	flag.Var(&f.Exclude, "exclude", "exclude a key from comparison (repeatable)")

	flag.Parse()
	return f
}
