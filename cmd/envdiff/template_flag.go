package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/envtemplate"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/loader"
)

var (
	tmplEnvFile    string
	tmplFile       string
	tmplMissingKey string
)

func registerTemplateFlags(fs *flag.FlagSet) {
	fs.StringVar(&tmplEnvFile, "tmpl-env", "", "env file providing substitution values")
	fs.StringVar(&tmplFile, "tmpl-file", "", "env file whose values are treated as templates")
	fs.StringVar(&tmplMissingKey, "tmpl-missing", "error", "missing key behaviour: error|zero|default")
}

func runTemplateCommand() {
	if tmplEnvFile == "" || tmplFile == "" {
		fmt.Fprintln(os.Stderr, "envdiff template: --tmpl-env and --tmpl-file are required")
		os.Exit(1)
	}

	envMap, err := loader.LoadFile(tmplEnvFile, filter.ApplyToMap, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff template: load env file: %v\n", err)
		os.Exit(1)
	}

	tmplMap, err := loader.LoadFile(tmplFile, filter.ApplyToMap, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff template: load template file: %v\n", err)
		os.Exit(1)
	}

	opts := envtemplate.Options{MissingKey: tmplMissingKey}
	out, err := envtemplate.Apply(tmplMap, envMap, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff template: render error: %v\n", err)
		os.Exit(1)
	}

	for k, v := range out {
		fmt.Printf("%s=%s\n", k, v)
	}
}
