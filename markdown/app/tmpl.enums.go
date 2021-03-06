package app

//TmplEnumsHandler 服务处理函数
const TmplEnumsHandlerDDS = `
package {{.PKG}}

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/types"
	"gitlab.100bm.cn/micro-plat/dds/dds"
)

//SystemEnumsHandler 枚举数据查询服务
type SystemEnumsHandler struct {
}

//NewSystemEnumsHandler 枚举数据查询服务
func NewSystemEnumsHandler() *SystemEnumsHandler {
	return &SystemEnumsHandler{}
}

//QueryHandle 枚举数据查询服务
func (o *SystemEnumsHandler) QueryHandle(ctx hydra.IContext) interface{} {

	//根据传入的枚举类型获取数据
	tp := ctx.Request().GetString("dic_type")
	if tp != "" {
		list, err := getEnums(tp, ctx.Request().GetMap())
		if err != nil {
			return err
		}
		return list
	}

	//查询所有枚举数据
	list := types.XMaps{}
	for _, sql := range enumsMap {
		items, err := hydra.C.DB().GetRegularDB().Query(sql, ctx.Request().GetMap())
		if err != nil {
			return err
		}
		list = append(list, items...)
	}
	return list
}

func getEnums(tp string, m types.XMap) (types.XMaps, error) {
	if _, ok := enumsMap[tp]; ok {
		items, err := hydra.C.DB().GetRegularDB().Query(enumsMap[tp], m)
		return items, err
	}

	if _, ok := enumsExt[tp]; ok {
		items, err := hydra.C.DB().GetRegularDB().Query(enumsExt[tp], m)
		return items, err
	}

	items, err := dds.GetEnums(nil, m)
	return items, err
}

var enumsExt = map[string]string{}

var enumsMap = map[string]string{
{{ range $j,$t:=.Tbs -}}
{{if $t|fIsEnumTB -}}
{{$count:= 0 -}}
"{{$t.Name|rmhd|lower}}":{###}select {{if not ($t|fHasDT) -}} '{{$t.Name|rmhd}}' type{{$count = 1}}{{end -}}
{{- range $i,$c:=.Rows -}}
{{if $c.Con|fIsDI -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} value{{end -}}
{{if $c.Con|fIsDN -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} name{{end -}}
{{if $c.Con|fIsDT -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} type{{end -}}
{{if $c.Con|fIsDPID -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} pid{{end -}}
{{if $c.Con|fIsDC -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}}{{end -}}
{{end}} from {{$t.Name}} t {###},
{{end -}}
{{- end -}}
}`

//TmplEnumsHandler 服务处理函数
const TmplEnumsHandler = `
package {{.PKG}}

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/types"
)

//SystemEnumsHandler 枚举数据查询服务
type SystemEnumsHandler struct {
}

//NewSystemEnumsHandler 枚举数据查询服务
func NewSystemEnumsHandler() *SystemEnumsHandler {
	return &SystemEnumsHandler{}
}

//QueryHandle 枚举数据查询服务
func (o *SystemEnumsHandler) QueryHandle(ctx hydra.IContext) interface{} {

	//根据传入的枚举类型获取数据
	tp := ctx.Request().GetString("dic_type")
	if tp != "" {
		list, err := getEnums(tp, ctx.Request().GetMap())
		if err != nil {
			return err
		}
		return list
	}

	//查询所有枚举数据
	list := types.XMaps{}
	for _, sql := range enumsMap {
		items, err := hydra.C.DB().GetRegularDB().Query(sql, ctx.Request().GetMap())
		if err != nil {
			return err
		}
		list = append(list, items...)
	}
	return list
}

func getEnums(tp string, m types.XMap) (types.XMaps, error) {
	if _, ok := enumsMap[tp]; ok {
		items, err := hydra.C.DB().GetRegularDB().Query(enumsMap[tp], m)
		return items, err
	}

	if _, ok := enumsExt[tp]; ok {
		items, err := hydra.C.DB().GetRegularDB().Query(enumsExt[tp], m)
		return items, err
	}

	items, err := dds.GetEnums(nil, m)
	return items, err
}

var enumsExt = map[string]string{}

var enumsMap = map[string]string{
{{ range $j,$t:=.Tbs -}}
{{if $t|fIsEnumTB -}}
{{$count:= 0 -}}
"{{$t.Name|rmhd|lower}}":{###}select {{if not ($t|fHasDT) -}} '{{$t.Name|rmhd}}' type{{$count = 1}}{{end -}}
{{- range $i,$c:=.Rows -}}
{{if $c.Con|fIsDI -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} value{{end -}}
{{if $c.Con|fIsDN -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} name{{end -}}
{{if $c.Con|fIsDT -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} type{{end -}}
{{if $c.Con|fIsDC -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}}{{end -}}
{{end}} from {{$t.Name}} t {###},
{{end -}}
{{- end -}}
}`
