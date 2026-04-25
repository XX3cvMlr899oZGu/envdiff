// Package envmatrix provides a multi-environment comparison matrix for env maps.
//
// Given a set of named environments (e.g. dev, staging, prod), Build produces
// a slice of Row values — one per unique key — each carrying the value seen in
// every environment and a list of environments where the key is absent.
//
// Example:
//
//	envs := map[string]map[string]string{
//		"dev":  {"PORT": "8080", "DEBUG": "true"},
//		"prod": {"PORT": "443",  "DEBUG": "false"},
//	}
//	rows := envmatrix.Build(envs, envmatrix.DefaultOptions())
//	fmt.Print(envmatrix.FormatText(rows, []string{"dev", "prod"}))
package envmatrix
