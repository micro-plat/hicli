package tmpl

const MdDictionaryTPL = `
{{ $empty := "" -}}
###  {{.Desc}}[{{.Name}}]

| 字段名       | 类型       | 默认值  | 为空 |   约束    | 描述                                |
| ------------| -----------| :----: | :--: | :-------: | :---------------------------------|
{{range $i,$c:=.Rows -}}
| {{$c.Name}} | {{$c.Type}}{{if ne $c.LenStr $empty}}({{$c.LenStr}}){{end}}|{{$c.Def}}|{{$c|IsMDNull}}| {{$c.Con}} | {{$c.Desc}}|
{{end -}}
`