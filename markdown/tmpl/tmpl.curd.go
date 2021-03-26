// +build !oracle

package tmpl

const MarkdownCurdSql = `
{{- $length := 32 -}}
{{- $createrows := .Rows|create -}}
{{- $updaterows := .Rows|update -}}
{{- $detailrows := .Rows|detail -}}
{{- $deleterows := .Rows|delete -}}
{{- $listrows := .Rows|list -}}
{{- $queryrows := .Rows|query -}}
{{- $pks := .|pks -}}
{{- $order:=.Rows|order|orderSort -}}
{{- $sort:=.Rows|sort -}}
{{- $empty:="" -}}
package {{.PKG}}

{{- if (gt ($createrows|len) 0)}}
//Insert{{.Name|rmhd|upperName}} 添加{{.Desc}}
const Insert{{.Name|rmhd|upperName}} = {###}
insert into {{.Name}}{{.DBLink}}
(
	{{if (.|seq) }}{{range $i,$c:=$pks}}{{$c}},{{end}}{{end}}
	{{- range $i,$c:=$createrows}}
	{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}
	{{- end}}
)
values
(
	{{if (.|seq)}}{{range $i,$c:=$pks}}@{{$c}},{{end}}{{end}}
	{{- range $i,$c:=$createrows}}
	{{if or ($c.Type|codeType|isInt) ($c.Type|codeType|isInt64) ($c.Type|codeType|isDecimal) }}if(isnull(@{{$c.Name}})||@{{$c.Name}}='',0,@{{$c.Name}}){{if lt $i ($createrows|maxIndex)}},{{end}}{{else -}}
	@{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}{{end}}
	{{- end}}
){###}
{{end -}}

{{if gt ($detailrows|len) 0 -}}
//Get{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} 查询{{.Desc}}单条数据
const Get{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} = {###}
select 
{{- range $i,$c:=$detailrows}}
	t.{{$c.Name}}{{if lt $i ($detailrows|maxIndex)}},{{end}}
{{- end}}
from {{.Name}} t
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&{{$c}} 
{{- end}}{{end}}{###}

{{- if gt (.TabInfo.TabField|len) 0}}
//Get{{.Name|rmhd|upperName}}Detail 查询{{.Desc}}单条详情数据
const Get{{.Name|rmhd|upperName}}Detail= {###}
select 
{{- range $i,$c:=$detailrows}}
	t.{{$c.Name}}{{if lt $i ($detailrows|maxIndex)}},{{end}}
{{- end}}
from {{.Name}} t
where
	{{- range $i,$c:=.TabInfo.TabField}}
	&{{(or ($c) ($pks|firstStr))}}
	{{- end}}
{###}
{{- end}}
{{- end}}

//Get{{.Name|rmhd|upperName}}ListCount 获取{{.Desc}}列表条数
const Get{{.Name|rmhd|upperName}}ListCount = {###}
select count(1)
from {{.Name}} t
where 
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$deleterows}}
	and {{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=$queryrows -}}
{{if $c.Type|codeType|isTime }}
	and t.{{$c.Name}} >= @{{$c.Name}} 
	and t.{{$c.Name}} < date_add(@{{$c.Name}}, interval 1 day)
{{- else if and ($c.Type|codeType|isString) (gt $c.Len $length)}}
	?t.{{$c.Name}}
{{- else}}
	&t.{{$c.Name}}{{end}}
{{- end}}{{end}}{###}

//Get{{.Name|rmhd|upperName}}List 查询{{.Desc}}列表数据
const Get{{.Name|rmhd|upperName}}List = {###}
select 
{{- range $i,$c:=$listrows}}
	t.{{$c.Name}}{{if lt $i ($listrows|maxIndex)}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$deleterows}}
	and {{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=$queryrows -}}
{{if $c.Type|codeType|isTime }}
	and t.{{$c.Name}} >= @{{$c.Name}} 
	and t.{{$c.Name}} < date_add(@{{$c.Name}}, interval 1 day)
{{- else if and ($c.Type|codeType|isString)  (gt $c.Len $length)}}
	?t.{{$c.Name}}
{{- else}}
	&t.{{$c.Name}}{{end}}
{{- end}} 
{{- if gt ($sort|len) 0}}
order by #order_by
{{- else if gt ($order|len) 0}}
order by {{range $i,$c:=$order}}t.{{$c.Name}} {{or ($c.Con|orderCon) "desc"}}{{if lt $i ($order|maxIndex)}}, {{end}}{{end}}
{{- else}}
order by {{range $i,$c:=$pks}}t.{{$c}} desc{{end}}
{{- end}}
limit @ps offset @offset
{{end -}}{###}

{{- if gt (.TabInfo.TabListField|len) 0}}
//Get{{.Name|rmhd|upperName}}DetailListCount 获取{{.Desc}}列表条数
const Get{{.Name|rmhd|upperName}}DetailListCount = {###}
select count(1)
from {{.Name}} t
where 
{{- range $i,$c:=$deleterows}}
	and {{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=.TabInfo.TabListField}}
&{{(or ($c) ($pks|firstStr))}}
{{- end}}{###}

//Get{{.Name|rmhd|upperName}}DetailList 查询{{.Desc}}列表数据
const Get{{.Name|rmhd|upperName}}DetailList = {###}
select 
{{- range $i,$c:=$listrows}}
	t.{{$c.Name}}{{if lt $i ($listrows|maxIndex)}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- range $i,$c:=$deleterows}}
	and {{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=.TabInfo.TabListField}}
&{{(or ($c) ($pks|firstStr))}}
{{- end}}
limit @ps offset @offset{###}
{{- end}}
{{end}}


{{- if (gt ($updaterows|len) 0)}}
//Update{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} 更新{{.Desc}}
const Update{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} = {###}
update {{.Name}}{{.DBLink}} 
set
{{- range $i,$c:=$updaterows}}
	{{$c.Name}} =	{{if or ($c.Type|codeType|isInt) ($c.Type|codeType|isInt64) ($c.Type|codeType|isDecimal) }}if(isnull(@{{$c.Name}})||@{{$c.Name}}='',0,@{{$c.Name}}){{if lt $i ($updaterows|maxIndex)}},{{end}}{{else -}}
	@{{$c.Name}}{{if lt $i ($updaterows|maxIndex)}},{{end}}{{end}}
{{- end}}
where
{{- if eq ($pks|len) 0}}
	1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&{{$c}}
{{- end}}{{end}}{###}
{{end -}}

{{- if gt ($deleterows|len) 0}}
//Delete{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} 删除{{.Desc}}
const Delete{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} = {###}
update {{.Name}}{{.DBLink}} 
set
{{- range $i,$c:=$deleterows}}
	{{$c.Name}}={{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&{{$c}}
{{- end}}
{{- range $i,$c:=$deleterows}}
	and {{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{end }}{###}
{{end}}
`
