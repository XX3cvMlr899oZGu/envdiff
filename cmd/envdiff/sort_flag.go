package main

import (
	"flag"

	"github.com/user/envdiff/internal/envsort"
)

var (
	sortOrder  string
	sortPrefix string
)

func registerSortFlags() {
	flag.StringVar(&sortOrder, "sort", "asc", "Sort order for output keys: asc or desc")
	flag.StringVar(&sortPrefix, "sort-prefix", "", "Keys with this prefix are listed first")
}

func buildSortOptions() envsort.Options {
	order := envsort.Ascending
	if sortOrder == "desc" {
		order = envsort.Descending
	}
	return envsort.Options{
		Order:  order,
		Prefix: sortPrefix,
	}
}
