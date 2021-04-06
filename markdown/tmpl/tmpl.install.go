package tmpl

const InstallTmpl = `
package {{.PKG}}

import (
	"github.com/micro-plat/hydra"
	_ "github.com/go-sql-driver/mysql"
)
		
func init() {
	//注册服务包
	hydra.OnReadying(func() error {
		hydra.Installer.DB.AddSQL(
		{{range $i,$c:=.Tbs -}}
		{{$c.Name}},
		{{end -}}
		{{if .SEQFile}}SEQ_IDS,{{end}}
		)
		return nil
	}) 
}
`
