package tmpl

const UpdateSingle = `
{{$count:=.Rows|maxIndex -}}
{{$rcount:=.|pks|maxIndex -}}
//Update{{.Name|rmhd|varName}} 更新{{.Desc}}
const Update{{.Name|rmhd|varName}} = {###}
update {{.Name}} t set
{{- range $i,$c:=.Rows}}
t.{{$c.Name}} = @{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
where
{{- range $i,$c:=.|pks}}
t.{{$c}} = @{{$c}}{{if lt $i $rcount}} and {{end}}
{{- end}} 
{###}`
