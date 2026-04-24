// Package envtemplate renders Go text/template expressions embedded inside
// .env values, substituting variables from a provided environment map.
//
// Example:
//
//	env := map[string]string{"DB_USER": "admin", "DB_PASS": "s3cr3t"}
//	tmpl := map[string]string{"DSN": "postgres://{{.DB_USER}}:{{.DB_PASS}}@localhost/app"}
//	out, err := envtemplate.Apply(tmpl, env, envtemplate.DefaultOptions())
//	// out["DSN"] == "postgres://admin:s3cr3t@localhost/app"
//
// Missing key behaviour is controlled by Options.MissingKey, which maps
// directly to Go's text/template missingkey option ("error", "zero", "default").
package envtemplate
