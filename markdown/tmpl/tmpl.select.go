package tmpl

const SelectSingle = `
{{$count:=.Rows|maxIndex -}}
{{$rcount:=.|pks|maxIndex -}}
//Select{{.Name|rmhd|varName}} 查询{{.Desc}}
const Select{{.Name|rmhd|varName}} = {###}
select 
{{- range $i,$c:=.Rows}}
t.{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- range $i,$c:=.|pks}}
t.{{$c}} = @{{$c}}{{if lt $i $rcount}} and {{end}}
{{- end}} 
{###}`
