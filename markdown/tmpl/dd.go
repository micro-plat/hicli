package tmpl

const MdOracleTPL = `
{{ $empty := "" -}}
###  {{.desc}}[{{.name}}]

| 字段名       | 类型       | 默认值  | 为空 |   约束    | 描述                                |
| ------------| -----------| :----: | :--: | :-------: | :---------------------------------|
{{range $i,$c:=.columns -}}
| {{$c.name}} | {{$c.type}}{{if ne $c.lenstr $empty}}({{$c.lenstr}}){{end}}|{{$c.def}}|{{if $c.isnull}}是{{else}}否{{end}}| {{$c.cons}} | {{$c.desc}}|
{{end -}}
`

const MdMysqlTPL = `
{{ $empty := "" -}}
###  {{.Desc}}[{{.Name}}]

| 字段名       | 类型       | 默认值  | 为空 |   约束    | 描述                                |
| ------------| -----------| :----: | :--: | :-------: | :---------------------------------|
{{range $i,$c:=.Rows -}}
| {{$c.Name}} | {{$c.Type}}{{if ne $c.LenStr $empty}}({{$c.LenStr}}){{end}}|{{$c.Def}}|{{$c|IsMDNull}}| {{$c.Con}} | {{$c.Desc}}|
{{end -}}
`
