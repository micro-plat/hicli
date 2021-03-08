package tmpl

const MarkdownCurdDriverSql = `
package sql

import (
	_ "github.com/go-sql-driver/mysql"
	_ "{{.BasePath}}/modules/const/sql/mysql"
)
`
