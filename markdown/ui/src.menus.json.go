package ui

const SrcMenusJson = `
{{- $rows:=. -}}
[
{{- range $i,$v:=$rows}}
  {
    "name": "{{$v.Desc}}",
    "path": "/{{$v.Name|rmhd|rpath}}"
  }{{if lt $i ($rows|maxIndex)}},{{end}}
{{- end}}
]
`
