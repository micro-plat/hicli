//go:build oracle
// +build oracle

package tmpl

const MarkdownCurdDriverSql = `
package sql

import (
	_ "github.com/mattn/go-oci8"
)
`
