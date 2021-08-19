package ui

const SrcMenusJson = `
{{- $rows:=.router -}}
{{- $ext:=.ext -}}
[
{{- range $i,$v:=$rows}}
  {
    "name": "{{$v.Desc}}",
    "path": "/{{$v.Name|rmhd|rpath}}"
  }{{if or (lt $i ($rows|maxIndex)) (gt ($ext|len) 0)}},{{end}}
{{- end}}
{{- range $i,$v:=$ext}}
  {
    "name": "{{$v.Desc}}",
    "path": "/{{$v.Path}}"
  }{{if lt $i ($ext|maxIndex)}},{{end}}
{{- end}}
]
`
