package tmpl

const FieldsTmpl = `package field

{{range $j,$r:=.}}
//Field{{$r.Name|varName}} 字段{{.Desc}}的数据库名称
const Field{{$r.Name|varName}} = "{{$r.Name}}"
{{end -}}
`
