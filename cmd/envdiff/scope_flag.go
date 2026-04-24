package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/envscope"
)

var (
	scopePrefix      string
	scopeName        string
	scopeStripPrefix bool
	scopeCaseFold    bool
)

func registerScopeFlags() {
	flag.StringVar(&scopePrefix, "scope-prefix", "", "key prefix to scope (e.g. APP_)")
	flag.StringVar(&scopeName, "scope-name", "default", "name for the scope")
	flag.BoolVar(&scopeStripPrefix, "scope-strip", false, "strip prefix from scoped keys")
	flag.BoolVar(&scopeCaseFold, "scope-casefold", false, "case-insensitive prefix matching")
}

func buildScopeOptions() (envscope.Scope, envscope.Options) {
	scope := envscope.Scope{
		Name:   strings.TrimSpace(scopeName),
		Prefix: strings.TrimSpace(scopePrefix),
	}
	opts := envscope.Options{
		StripPrefix: scopeStripPrefix,
		CaseFold:    scopeCaseFold,
	}
	return scope, opts
}

func printScopedEnv(scoped map[string]string) {
	for k, v := range scoped {
		fmt.Printf("%s=%s\n", k, v)
	}
}
