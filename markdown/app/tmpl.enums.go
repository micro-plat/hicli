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
		var items types.XMaps
		var err error
		if _, ok := enumsMap[tp]; !ok {
			items, err = dds.GetEnums(ctx, ctx.Request().GetMap())
		} else {
			items, err = hydra.C.DB().GetRegularDB().Query(enumsMap[tp], ctx.Request().GetMap())
		}
		if err != nil {
			return err
		}
		return items
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

var enumsMap = map[string]string{
{{ range $j,$t:=.Tbs -}}
{{if $t|fIsEnumTB -}}
{{$count:= 0 -}}
"{{$t.Name|rmhd|lower}}":{###}select {{if not ($t|fHasDT) -}} '{{$t.Name|rmhd}}' type {{$count = 1}}{{end -}}
{{- range $i,$c:=.Rows -}}
{{if $c.Con|fIsDI -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} value {{end -}}
{{if $c.Con|fIsDN -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} name {{end -}}
{{if $c.Con|fIsDT -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} type {{end -}}
{{end}} from {{$t.Name}} t {{if gt (.Rows|delete|len) 0}}where{{end}}
{{- range $i,$c:=.Rows|delete}}	and t.{{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($t.Rows|delete|maxIndex)}},{{end}}{{- end}}{###},
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
		if _, ok := enumsMap[tp]; ok {
			items, err := hydra.C.DB().GetRegularDB().Query(enumsMap[tp], ctx.Request().GetMap())
			if err != nil {
				return err
			}
			return items
		}
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

var enumsMap = map[string]string{
{{ range $j,$t:=.Tbs -}}
{{if $t|fIsEnumTB -}}
{{$count:= 0 -}}
"{{$t.Name|rmhd|lower}}":{###}select {{if not ($t|fHasDT) -}} '{{$t.Name|rmhd}}' type {{$count = 1}}{{end -}}
{{- range $i,$c:=.Rows -}}
{{if $c.Con|fIsDI -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} value {{end -}}
{{if $c.Con|fIsDN -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} name {{end -}}
{{if $c.Con|fIsDT -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} type {{end -}}
{{end}} from {{$t.Name}} t {{if gt (.Rows|delete|len) 0}}where{{end}}
{{- range $i,$c:=.Rows|delete}}	and t.{{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($t.Rows|delete|maxIndex)}},{{end}}{{- end}}{###},
{{end -}}
{{- end -}}
}`
