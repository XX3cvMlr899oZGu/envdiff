// Package envgroup partitions a flat env map into named groups based on
// key prefixes. For example, keys DB_HOST and DB_PORT are placed into a
// group named "DB", while APP_NAME goes into "APP".
//
// This is useful for inspecting large .env files where keys are
// conventionally namespaced by service or component.
//
// Usage:
//
//	groups := envgroup.Apply(env, envgroup.DefaultOptions())
//	for _, g := range groups {
//	    fmt.Println(g.Name, g.Keys)
//	}
package envgroup
