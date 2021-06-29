package tmpl

const MdDBDataTPL = `
{{- if .PKG}}package {{.PKG}}{{- end}}

import "github.com/micro-plat/hydra"

func init() {
	hydra.OnReady(func() error {
		hydra.Installer.DB.AddSQL(
			DataList...,
		)
		return nil
	})
}

//DataList 初始需导入数据
var DataList = []string{
{{.DataList}}
}
`
