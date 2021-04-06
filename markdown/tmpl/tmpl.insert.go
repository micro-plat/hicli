package tmpl

const InsertSingle = `
{{$count:=.Rows|maxIndex -}}
{{$rcount:=.|pks|maxIndex -}}
//Insert{{.Name|rmhd|varName}} 插入{{.Desc}}
const Insert{{.Name|rmhd|varName}} = {###}
insert into {{.Name}}(
{{- range $i,$c:=.Rows}}
	{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
)values(
{{- range $i,$c:=.Rows}}
	@{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}}
)
{###}`
