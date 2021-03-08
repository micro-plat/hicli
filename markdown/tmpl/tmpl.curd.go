package tmpl

const MarkdownCurdSql = `
{{- $length := 32 -}}
{{- $createrows := .Rows|create -}}
{{- $updaterows := .Rows|update -}}
{{- $detailrows := .Rows|detail -}}
{{- $deleterows := .Rows|delete -}}
{{- $listrows := .Rows|list -}}
{{- $queryrows := .Rows|query -}}
{{- $ismysql := .DBType|ismysql -}}
{{- $isoracle := .DBType|isoracle -}}
{{- $pks := .|pks -}}
{{- $order:=.|order -}}
{{- $empty:="" -}}
package {{.PKG}}

{{- if and $ismysql (gt ($createrows|len) 0)}}
//Insert{{.Name|rmhd|upperName}} 添加{{.Desc}}
const Insert{{.Name|rmhd|upperName}} = {###}
insert into {{.Name}}{{.DBLink}}
(
	{{if ne (.|seqValue) $empty}}{{range $i,$c:=$pks}}{{$c}},{{end}}{{end}}
	{{- range $i,$c:=$createrows}}
	{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}
	{{- end}}
)
values
(
	{{if ne (.|seqValue) $empty}}{{range $i,$c:=$pks}}@{{$c}},{{end}}{{end}}
	{{- range $i,$c:=$createrows}}
	{{if or ($c.Type|codeType|isInt) ($c.Type|codeType|isInt64) ($c.Type|codeType|isDecimal) }}if(isnull(@{{$c.Name}})||@{{$c.Name}}='',0,@{{$c.Name}}){{if lt $i ($createrows|maxIndex)}},{{end}}{{else -}}
	@{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}{{end}}
	{{- end}}
){###}
{{end -}}

{{- if and $isoracle (gt ($createrows|len) 0)}}
//Insert{{.Name|rmhd|upperName}} 添加{{.Desc}}
const Insert{{.Name|rmhd|upperName}} = {###}
insert into {{.Name}}{{.DBLink}}
(
	{{- range $i,$c:=.seqs}}
	{{$c.Name}},{{end}}{{range $i,$c:=$createrows}}
	{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}
	{{- end}}
)
values(
	{{- range $i,$c:=.seqs}}
	{{$c.seqname}}.nextval,
	{{- end}}
	{{- range $i,$c:=$createrows}}
	{{if $c.Type|codeType|isTime }}to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss'){{else -}}
	@{{$c.Name}}{{end}}{{if $c.comma}},{{end}}
	{{- end}}
){###}
{{end -}}

{{- if $ismysql}}
{{if gt ($detailrows|len) 0 -}}
//Get{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} 查询单条数据{{.Desc}}
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
{{- end}}

//Get{{.Name|rmhd|upperName}}ListCount 获取{{.Desc}}列表条数
const Get{{.Name|rmhd|upperName}}ListCount = {###}
select count(1)
from {{.Name}} t
where 
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$queryrows -}}
{{if $c.Type|codeType|isTime }}
	and t.{{$c.Name}} >= @{{$c.Name}} 
	and t.{{$c.Name}} < date_add(@{{$c.Name}}, interval 1 day)
{{- else if ($c.Type|codeType|isString)}}
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
{{- range $i,$c:=$queryrows -}}
{{if $c.Type|codeType|isTime }}
	and t.{{$c.Name}} >= @{{$c.Name}} 
	and t.{{$c.Name}} < date_add(@{{$c.Name}}, interval 1 day)
{{- else if ($c.Type|codeType|isString)}}
	?t.{{$c.Name}}
{{- else}}
	&t.{{$c.Name}}{{end}}
{{- end}} 
{{- if gt ($order|len) 0}}
order by {{range $i,$c:=$order}}t.{{$c.Name}}{{if $c.comma}},{{else}} desc{{end}}{{end}}
{{- else}}
order by {{range $i,$c:=$pks}}t.{{$c}} desc{{end}}
{{- end}}
limit @ps offset @offset
{{end -}}{###}{{end}}

{{- if and $isoracle }}
{{if  (gt ($detailrows|len) 0) -}}
//Get{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} 查询单条数据{{.Desc}}
const Get{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} = {###}
select 
{{- range $i,$c:=$detailrows}}
	t.{{$c.Name}}{{if lt $i ($detailrows|maxIndex)}},{{end}}
{{- end}} 
from {{.Name}}{{.DBLink}} t
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&{{$c}}
{{- end}}{{- end}}{###}
{{- end}}

//Get{{.Name|rmhd|upperName}}ListCount 获取{{.Desc}}列表条数
const Get{{.Name|rmhd|upperName}}ListCount = {###}
select count(1)
from {{.Name}}{{.DBLink}} t
where
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$queryrows -}}
{{if $c.Type|codeType|isTime }}
	and t.{{$c.Name}} >= to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')
  and t.{{$c.Name}} < to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')+1
{{- else if ($c.Type|codeType|isString)}}
  ?t.{{$c.Name}}
{{- else}}
	&t.{{$c.Name}}{{end}}
{{- end}}{{end}}
{###}

//Get{{.Name|rmhd|upperName}}List 查询{{.Desc}}列表数据
const Get{{.Name|rmhd|upperName}}List = {###}
select 
	TAB1.*
from (select L.*  
	from (select rownum as rn,R.* 
		from (
			select 
			{{- range $i,$c:=$listrows}}
			{{if $c.Type|codeType|isTime }}to_char(t.{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')	{{$c.Name}}{{else -}}
				t.{{$c.Name}}{{end}}{{if lt $i ($listrows|maxIndex)}},{{end}}
			{{- end}} 
			from {{.Name}}{{.DBLink}} t
			where
			{{- if eq ($listrows|len) 0}}
				1=1
			{{- else -}}
			{{- range $i,$c:=$queryrows -}} 
			{{if $c.Type|codeType|isTime }}
				and t.{{$c.Name}} >= to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')
				and t.{{$c.Name}} < to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')+1
			{{- else if ($c.Type|codeType|isString)}}
				?t.{{$c.Name}}
			{{- else}}
				&t.{{$c.Name}}{{end}}
			{{- end}}{{end}}
			{{- if gt ($order|len) 0}}
			order by {{range $i,$c:=$order}}t.{{$c.Name}}{{if $c.comma}},{{else}} desc{{end}}{{end}}
			{{- else}}
			order by {{range $i,$c:=pks}}t.{{$c}} desc{{end}}
			{{- end}}
			) R 
	where rownum <= @pi * @ps) L 
where L.rn > (@pi - 1) * @ps) TAB1{###}
{{end}}

{{- if  and $ismysql (gt ($updaterows|len) 0)}}
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

{{- if and $isoracle (gt ($updaterows|len) 0)}}
//Update{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} 更新{{.Desc}}
const Update{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}} = {###}
update {{.Name}}{{.DBLink}} 
set
{{- range $i,$c:=$updaterows}}
	{{if $c.Type|codeType|isTime }}{{$c.Name}}=to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss'){{else -}}
	{{$c.Name}} = @{{$c.Name}}{{end}}{{if lt $i ($updaterows|maxIndex)}},{{end}}
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
