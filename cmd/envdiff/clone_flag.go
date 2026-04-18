package main

import (
	"flag"
	"strings"

	"github.com/user/envdiff/internal/envclone"
)

var (
	clonePrefix      string
	cloneStripPrefix bool
	cloneLowercase   bool
)

func registerCloneFlags() {
	flag.StringVar(&clonePrefix, "clone-prefix", "", "filter keys by prefix when cloning")
	flag.BoolVar(&cloneStripPrefix, "clone-strip-prefix", false, "strip prefix from cloned keys")
	flag.BoolVar(&cloneLowercase, "clone-lowercase", false, "lowercase all keys in cloned output")
}

func buildCloneOptions() envclone.Options {
	opts := envclone.Options{
		KeyPrefix:   clonePrefix,
		StripPrefix: cloneStripPrefix,
	}
	if cloneLowercase {
		opts.KeyTransform = strings.ToLower
	}
	return opts
}
