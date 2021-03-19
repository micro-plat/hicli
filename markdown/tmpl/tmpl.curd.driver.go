package tmpl

const MarkdownCurdDriverSql = `
{{- if .IsOracle -}}
// +build oracle
{{else}}
// +build mysql
{{end}}
package sql

import (
	{{- if not .IsOracle}}
	_ "github.com/go-sql-driver/mysql"
	_ "{{.BasePath}}/modules/const/sql/mysql"
	{{- else}}
	_ "github.com/mattn/go-oci8"
  {{- end }}
)
`
