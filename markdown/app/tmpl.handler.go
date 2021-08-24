// +build !oracle

package app

//TmplServiceHandler 服务处理函数
const TmplServiceHandler = `
{{- $empty := "" -}}
{{- $rows := .Rows -}}
{{- $pks := .|pks -}}
{{- $sort:=.Rows|sort -}}
{{- $btns:=.BtnInfo -}}
{{- $up:= 0 -}}
{{- range $i,$c:=$rows|update}}{{if $c.Con|UP}}{{$up = 1}}{{end}}{{end -}}
{{- $db:= false -}}
{{- range $i,$c:=.Rows|query}}{{- if ($c.Con|CSCR) }}{{$db = true}}{{end}}{{end -}}
package {{.PKG}}

import (
	{{- if eq $up 1}}
	"fmt"
	"io"
	"os"
	"path"
	"time"
	"path/filepath"
	{{- end}}
	{{- if or (gt ($rows|list|len) 0) (gt ($rows|create|len) 0) (gt ($rows|detail|len) 0) (gt ($rows|update|len) 0) (gt ($rows|delete|len) 0)}}
	"net/http"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
	"{{.BasePath}}/modules/const/sql"
	"{{.BasePath}}/modules/const/field"
	{{- end}}
	{{- if gt ($rows|list|len) 0}}
	"github.com/micro-plat/lib4go/types"
	{{- end}}
	{{- if gt ($sort|len) 0}}
	"regexp"
	{{- end}}
	{{- if or (and (.|mysqlseq) (gt (.Rows|create|len) 0)) $db}}
	"{{.BasePath}}/modules/db"
	{{- end}}
)

//{{.Name|rmhd|varName}}Handler {{.Desc}}处理服务
type {{.Name|rmhd|varName}}Handler struct {
}

func New{{.Name|rmhd|varName}}Handler() *{{.Name|rmhd|varName}}Handler {
	return &{{.Name|rmhd|varName}}Handler{}
}

{{if gt (.Rows|create|len) 0 -}}
//PostHandle 添加{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) PostHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------添加{{.Desc}}数据--------")
	
	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(post{{.Name|rmhd|varName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	{{- if (.|mysqlseq) }}
	xdb := hydra.C.DB().GetRegularDB()
	{{$pks|firstStr|lowerName}}, err := db.GetNewID(xdb, sql.SQLGetSEQ, map[string]interface{}{"name": "{{.Desc}}"})
	if err != nil {
		return err
	}
	input := ctx.Request().GetMap()
	input["{{$pks|firstStr}}"] = {{$pks|firstStr|lowerName}}
	count, err := xdb.Execute(sql.Insert{{.Name|rmhd|upperName}}, input)
  {{- else}}
	count, err := hydra.C.DB().GetRegularDB().Execute(sql.Insert{{.Name|rmhd|upperName}}, ctx.Request().GetMap())
	{{- end}}
	if err != nil || count < 1 {
		return errs.NewErrorf(http.StatusNotExtended, "添加数据出错:%+v", err)
	}

	ctx.Log().Info("3.返回结果")
	return "success"
}
{{end}}


{{- if gt ($rows|detail|len) 0}}
//GetHandle 获取{{.Desc}}单条数据
func (u *{{.Name|rmhd|varName}}Handler) GetHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}单条数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(get{{.Name|rmhd|varName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	items, err := hydra.C.DB().GetRegularDB().Query(sql.Get{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}}, ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended, "查询数据出错:%+v", err)
	}
	if items.Len() == 0 {
		return errs.NewError(http.StatusNoContent, "未查询到数据")
	}

	ctx.Log().Info("3.返回结果")
	return items.Get(0)
}
{{if gt (.TabInfo.TabField|len) 0}}
//DetailHandle 获取{{.Desc}}详情单条数据
func (u *{{.Name|rmhd|varName}}Handler) DetailHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}详情单条数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(get{{.Name|rmhd|varName}}DetailCheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	items, err :=  hydra.C.DB().GetRegularDB().Query(sql.Get{{.Name|rmhd|upperName}}Detail,ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"查询数据出错:%+v", err)
	}
	if items.Len() == 0 {
		return errs.NewError(http.StatusNoContent, "未查询到数据")
	}

	ctx.Log().Info("3.返回结果")
	return items.Get(0)
}
{{- end}}
{{- end}}

{{if gt ($rows|list|len) 0 -}}
//QueryHandle  获取{{.Desc}}数据列表
func (u *{{.Name|rmhd|varName}}Handler) QueryHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}数据列表--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(query{{.Name|rmhd|varName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}
	{{- if gt ($sort|len) 0}}
	orderBy := ctx.Request().GetString("order_by")
	if len(orderBy) > 1 && !regexp.MustCompile("^t.[A-Za-z0-9_,.\\s]+ (asc|desc)$").MatchString(orderBy) {
		return errs.NewErrorf(http.StatusNotAcceptable, "排序参数校验错误!")
	}
	{{- end}}

	ctx.Log().Info("2.执行操作")
	m := ctx.Request().GetMap()
	m["offset"] = (ctx.Request().GetInt("pi") - 1) * ctx.Request().GetInt("ps")
	{{- range $i,$c:=.Rows|query}}
	{{- if ($c.Con|CSCR) }}
	m["{{$c.Name}}"] = db.GetInStr(m.GetString("{{$c.Name}}"), "and t.{{$c.Name}} in (%s)")
	{{- end}}
	{{- end}}
	count, err := hydra.C.DB().GetRegularDB().Scalar(sql.Get{{.Name|rmhd|upperName}}ListCount, m)
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended, "查询数据数量出错:%+v", err)
	}
	
	var items types.XMaps
	if types.GetInt(count) > 0 {
		items, err = hydra.C.DB().GetRegularDB().Query(sql.Get{{.Name|rmhd|upperName}}List, m)
		if err != nil {
			return errs.NewErrorf(http.StatusNotExtended, "查询数据出错:%+v", err)
		}
	}

	ctx.Log().Info("3.返回结果")
	return map[string]interface{}{
		"items": items,
		"count": types.GetInt(count),
	}
}
{{if gt (.TabInfo.TabListField|len) 0}}
//QueryDetailHandle  获取{{.Desc}}数据(详情)列表
func (u *{{.Name|rmhd|varName}}Handler) QueryDetailHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}数据(详情)列表--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(query{{.Name|rmhd|varName}}DetailCheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	m := ctx.Request().GetMap()
	m["offset"] = (ctx.Request().GetInt("pi") - 1) * ctx.Request().GetInt("ps")

	count, err := hydra.C.DB().GetRegularDB().Scalar(sql.Get{{.Name|rmhd|upperName}}DetailListCount, m)
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended, "查询数据数量出错:%+v", err)
	}
	
	var items types.XMaps
	if types.GetInt(count) > 0 {
		items, err = hydra.C.DB().GetRegularDB().Query(sql.Get{{.Name|rmhd|upperName}}DetailList, m)
		if err != nil {
			return errs.NewErrorf(http.StatusNotExtended, "查询数据出错:%+v", err)
		}
	}

	ctx.Log().Info("3.返回结果")
	return map[string]interface{}{
		"items": items,
		"count": types.GetInt(count),
	}
}
{{- end}}
{{- end}}

{{- if gt ($rows|export|len) 0}}
//ExportHandle  获取{{.Desc}}数据导出列表
func (u *{{.Name|rmhd|varName}}Handler) ExportHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}数据导出列表--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(export{{.Name|rmhd|varName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}
	{{- if gt ($sort|len) 0}}
	orderBy := ctx.Request().GetString("order_by")
	if len(orderBy) > 1 && !regexp.MustCompile("^t.[A-Za-z0-9_,.\\s]+ (asc|desc)$").MatchString(orderBy) {
		return errs.NewErrorf(http.StatusNotAcceptable, "排序参数校验错误!")
	}
	{{- end}}

	ctx.Log().Info("2.执行操作")
	m := ctx.Request().GetMap()
	m["offset"] = (ctx.Request().GetInt("pi") - 1) * ctx.Request().GetInt("ps")

	items, err := hydra.C.DB().GetRegularDB().Query(sql.Get{{.Name|rmhd|upperName}}ExportList, m)
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended, "查询数据出错:%+v", err)
	}
	
	ctx.Log().Info("3.返回结果")
	return map[string]interface{}{
		"items": items,
	}
}
{{- end}}

{{- if eq $up 1}}
//UploadHandle 上传文件
func (u *{{.Name|rmhd|varName}}Handler) UploadHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------上传文件--------")

	ctx.Log().Info("1. 读取文件内容")
	fileName, uf, _, err := ctx.Request().GetFile("file")
	if err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "读取文件错误:%+v", err)
	}
	defer uf.Close()

	ctx.Log().Info("2. 构建文件名")
	storeName := fmt.Sprintf("%s/%d%s", time.Now().Format("20060102"), hydra.C.UUID(), path.Ext(fileName))
	filePath := filepath.Join("../upload", storeName)
	if _, err = os.Stat(filepath.Dir(filePath)); err != nil {
		if !os.IsNotExist(err) {
			return errs.NewError(http.StatusInternalServerError, err)
		}
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return errs.NewError(http.StatusInternalServerError, err)
		}
	}

	nf, err := os.Create(filePath)
	if err != nil {
		return errs.NewError(http.StatusInternalServerError, err)
	}
	defer nf.Close()

	ctx.Log().Info("3. 保存文件:", storeName)
	if _, err = io.Copy(nf, uf); err != nil {
		return err
	}

	ctx.Log().Info("4. 返回结果")
	return map[string]interface{}{
		"file_name": storeName,
	}
}
{{- end}}

{{- if gt ($rows|update|len) 0}}
//GetUpdateHandle 获取{{.Desc}}更新的数据
func (u *{{.Name|rmhd|varName}}Handler) GetUpdateHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}更新的数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(getUpdate{{.Name|rmhd|varName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	items, err :=  hydra.C.DB().GetRegularDB().Query(sql.GetUpdate{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}},ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"查询数据出错:%+v", err)
	}
	if items.Len() == 0 {
		return errs.NewError(http.StatusNoContent, "未查询到数据")
	}

	ctx.Log().Info("3.返回结果")
	return items.Get(0)
}

//PutHandle 更新{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) PutHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------更新{{.Desc}}数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(update{{.Name|rmhd|varName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	_, err := hydra.C.DB().GetRegularDB().Execute(sql.Update{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}},ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"更新数据出错:%+v", err)
	}

	ctx.Log().Info("3.返回结果")
	return "success"
}
{{- end}}


{{- range $i,$btn:=$btns }}
{{- if  $btn.Show }}
//Get{{$btn.Name|upperName}}Handle 获取{{$.Desc}}单条数据
func (u *{{$.Name|rmhd|varName}}Handler) Get{{$btn.Name|upperName}}Handle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{$.Desc}}单条数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(get{{$.Name|rmhd|varName}}{{$btn.Name|upperName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	items, err :=  hydra.C.DB().GetRegularDB().Query(sql.Get{{$.Name|rmhd|upperName}}{{$btn.Name|upperName}}By{{$pks|firstStr|upperName}},ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"查询数据出错:%+v", err)
	}
	if items.Len() == 0 {
		return errs.NewError(http.StatusNoContent, "未查询到数据")
	}

	ctx.Log().Info("3.返回结果")
	return items.Get(0)
}
{{- end}}
{{- end}}

{{- if gt ($rows|delete|len) 0}}
//DeleteHandle 删除{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) DeleteHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------删除{{.Desc}}数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(delete{{.Name|rmhd|varName}}CheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	count,err := hydra.C.DB().GetRegularDB().Execute(sql.Delete{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}}, ctx.Request().GetMap())
	if err != nil||count<1 {
		return errs.NewErrorf(http.StatusNotExtended,"删除数据出错:%+v", err)
	}

	ctx.Log().Info("3.返回结果")
	return "success"
}
{{- end}}

{{- if gt (.Rows|create|len) 0}}
var post{{.Name|rmhd|varName}}CheckFields = map[string]interface{}{
	{{- range $i,$c:=.Rows|create}}
	{{- if ne ($c|isNull) $empty}}
	field.{{$c.Name|varName}}:"required",
	{{- end}}
	{{- end}}
}
{{end}}

{{- if gt (.Rows|detail|len) 0}}
var get{{.Name|rmhd|varName}}CheckFields = map[string]interface{}{
	{{range $i,$c:=$pks}}field.{{$c|varName}}:"required",{{end}}
}
{{if gt (.TabInfo.TabField|len) 0}}
var get{{.Name|rmhd|varName}}DetailCheckFields = map[string]interface{}{
	{{- range $i,$c:=.TabInfo.TabField}}
	field.{{(or ($c) ($pks|firstStr))|varName}}:"required",
	{{- end}}
}
{{end}}
{{- end}}

{{- if gt (.Rows|list|len) 0}}
var query{{.Name|rmhd|varName}}CheckFields = map[string]interface{}{
	{{- range $i,$c:=.Rows|query}}
	field.{{$c.Name|varName}}:"required",
	{{- end}}
}
{{- if gt (.TabInfo.TabListField|len) 0}}

var query{{.Name|rmhd|varName}}DetailCheckFields = map[string]interface{}{
	{{- range $i,$c:=.TabInfo.TabListField}}
	field.{{(or ($c) ($pks|firstStr))|varName}}:"required",
	{{- end}}
}
{{- end}}
{{- end}}

{{- if gt (.Rows|export|len) 0}}

var export{{.Name|rmhd|varName}}CheckFields = map[string]interface{}{
	{{- range $i,$c:=.Rows|query}}
	field.{{$c.Name|varName}}:"required",
	{{- end}}
}
{{- end}}

{{if gt (.Rows|update|len) 0 -}}
var update{{.Name|rmhd|varName}}CheckFields = map[string]interface{}{
	{{- range $i,$c:=.Rows|update}}
	{{- if ne ($c|isNull) $empty}}
	field.{{$c.Name|varName}}:"required",
	{{- end}}
	{{- end}}
}

var getUpdate{{.Name|rmhd|varName}}CheckFields = map[string]interface{}{
	{{range $i,$c:=$pks}}field.{{$c|varName}}:"required",{{end}}
}
{{- end}}

{{if gt (.Rows|delete|len) 0 -}}
var delete{{.Name|rmhd|varName}}CheckFields = map[string]interface{}{
	{{range $i,$c:=$pks}}field.{{$c|varName}}:"required",{{end}}
}
{{- end}}


{{- range $i,$btn:=$btns }}
{{- if  $btn.Show }}
var get{{$.Name|rmhd|varName}}{{$btn.Name|upperName}}CheckFields = map[string]interface{}{
	{{range $i,$c:=$pks}}field.{{$c|varName}}:"required",{{end}}
}
{{- end}}
{{- end}}
`
