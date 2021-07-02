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
{{- $btns:=.BtnInfo -}}
package {{.PKG}}

{{- if (gt ($createrows|len) 0)}}
//Insert{{.Name|rmhd|upperName}} 添加{{.Desc}}
const Insert{{.Name|rmhd|upperName}} = {###}
insert into {{.Name}}{{.DBLink}}
(
	{{if (.|mysqlseq) }}{{range $i,$c:=$pks}}{{$c}},{{end}}{{end}}
	{{- range $i,$c:=$createrows}}
	{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}
	{{- end}}
)
values
(
	{{if (.|mysqlseq)}}{{range $i,$c:=$pks}}@{{$c}},{{end}}{{end}}
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
	{{- if and ($c.Type|codeType|isString) ($c|replace) }}
	{{$c|replace}} {{$c.Name}}
	{{- else}}
	t.{{$c.Name}}
	{{- end}}
	{{- if lt $i ($detailrows|maxIndex)}},{{end}}
{{- end}}
from {{.Name}} t
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&t.{{$c}} 
{{- end}}{{end}}{###}

{{- if gt (.TabInfo.TabField|len) 0}}
//Get{{.Name|rmhd|upperName}}Detail 查询{{.Desc}}单条详情数据
const Get{{.Name|rmhd|upperName}}Detail= {###}
select 
{{- range $i,$c:=$detailrows}}
	{{- if and ($c.Type|codeType|isString) ($c|replace) }}
	{{$c|replace}} {{$c.Name}}
	{{- else}}
	t.{{$c.Name}}
	{{- end}}
	{{- if lt $i ($detailrows|maxIndex)}},{{end}}
{{- end}}
from {{.Name}} t
where
	{{- range $i,$c:=.TabInfo.TabField}}
	&t.{{(or ($c) ($pks|firstStr))}}
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
	and t.{{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=$queryrows -}}
{{if $c.Type|codeType|isTime }}
	{{- if ($c.Con|DRANGE)}}
	and t.{{$c.Name}} >= @start_time
	and t.{{$c.Name}} < date_add(@end_time, interval 1 day)
	{{- else}}
	and t.{{$c.Name}} >= @{{$c.Name}} 
	and t.{{$c.Name}} < date_add(@{{$c.Name}}, interval 1 day)
	{{- end}}
{{- else if and ($c.Type|codeType|isString) (gt $c.Len $length)}}
	?t.{{$c.Name}}
{{- else}}
	&t.{{$c.Name}}{{end}}
{{- end}}{{end}}{###}

//Get{{.Name|rmhd|upperName}}List 查询{{.Desc}}列表数据
const Get{{.Name|rmhd|upperName}}List = {###}
select 
{{- range $i,$c:=$listrows}}
	{{- if and ($c.Type|codeType|isString) ($c|replace) }}
	{{$c|replace}} {{$c.Name}}
	{{- else}}
	t.{{$c.Name}}
	{{- end}}
	{{- if lt $i ($listrows|maxIndex)}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$deleterows}}
	and t.{{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=$queryrows -}}
{{if $c.Type|codeType|isTime }}
{{- if ($c.Con|DRANGE)}}
	and t.{{$c.Name}} >= @start_time
	and t.{{$c.Name}} < date_add(@end_time, interval 1 day)
	{{- else}}
	and t.{{$c.Name}} >= @{{$c.Name}} 
	and t.{{$c.Name}} < date_add(@{{$c.Name}}, interval 1 day)
	{{- end}}
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
	and t.{{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=.TabInfo.TabListField}}
&t.{{(or ($c) ($pks|firstStr))}}
{{- end}}{###}

//Get{{.Name|rmhd|upperName}}DetailList 查询{{.Desc}}列表数据
const Get{{.Name|rmhd|upperName}}DetailList = {###}
select 
{{- range $i,$c:=$listrows}}
	{{- if and ($c.Type|codeType|isString) ($c|replace) }}
	{{$c|replace}} {{$c.Name}}
	{{- else}}
	t.{{$c.Name}}
	{{- end}}
	{{- if lt $i ($listrows|maxIndex)}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- range $i,$c:=$deleterows}}
	and t.{{$c.Name}}<>{{or ($c.Con|delCon) "1"}}{{if lt $i ($deleterows|maxIndex)}},{{end}}
{{- end}}
{{- range $i,$c:=.TabInfo.TabListField}}
&t.{{(or ($c) ($pks|firstStr))}}
{{- end}}
limit @ps offset @offset{###}
{{- end}}


{{- if (gt ($updaterows|len) 0)}}
//GetUpdate{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} 查询{{.Desc}}单条数据
const GetUpdate{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} = {###}
select 
{{- range $i,$c:=$updaterows}}
	{{- if and ($c.Type|codeType|isString) ($c|replace) }}
	{{$c|replace}} {{$c.Name}}
	{{- else}}
	t.{{$c.Name}}
	{{- end}}
	{{- if lt $i ($updaterows|maxIndex)}},{{end}}
{{- end}}
from {{.Name}} t
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&t.{{$c}} 
{{- end}}{{end}}{###}

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


{{- range $i,$btn:=$btns }}
//Update{{$.Name|rmhd|upperName}}{{$btn.Name|upperName}}By{{$pks|firstStr|upperName}} 更新数据
const Update{{$.Name|rmhd|upperName}}{{$btn.Name|upperName}}By{{$pks|firstStr|upperName}} = {###}
update {{$.Name}}{{$.DBLink}} 
set
{{- range $i,$c:=$btn.Rows}}
{{- if not $c.Disable}}
	{{- if $c.Type|codeType|isTime }}
	{{$c.Name}}=to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')
	{{- else}}
	{{$c.Name}} = @{{$c.Name}}{{end}}{{if ne $c.Name $btn.LastRowIndex}},{{end}}
{{- end}}
{{- end}}
where
{{- if eq ($pks|len) 0}}
	1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&{{$c}}
{{- end}}
{{- end}}{###}

{{- if $btn.Show }}
//Get{{$.Name|rmhd|upperName}}{{$btn.Name|upperName}}By{{$pks|firstStr|upperName}} 查询单条数据{{$.Desc}}
const Get{{$.Name|rmhd|upperName}}{{$btn.Name|upperName}}By{{$pks|firstStr|upperName}} = {###}
select 
{{- range $i,$c:=$btn.Rows}}
	{{or $c.SQLAliasName "t"}}.{{$c.Name}}{{if lt $i ($btn.Rows|maxIndex)}},{{end}}
{{- end}} 
from {{$.Name}}{{$.DBLink}} t
{{- range $i,$c:=$btn.Table}}
left join {{$c.Name}}{{$.DBLink}} t{{$i}} on t.{{index $btn.RelativeShelfFiled $c.Name}} = t{{$i}}.{{index $btn.RelativeFiled $c.Name}} 
{{- end}}
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&t.{{$c}}
{{- end}}
{{- end}}{###}
{{- end}}

{{- end}}
`
