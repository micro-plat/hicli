package tmpl

const FieldsTmpl = `package field

{{range $j,$r:=.}}
//{{$r.Name|varName}} 字段{{.Desc}}的数据库名称
const {{$r.Name|varName}} = "{{$r.Name}}"
{{end -}}
`
