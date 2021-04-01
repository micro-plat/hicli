// +build !oracle

package tmpl

const SQLTmpl = `

{{- if .PKG}}package {{.PKG}}
{{end -}}

{{$count:=.Rows|maxIndex -}}

{{- if .PKG}} 
//{{.Name}} {{.Desc}}
const {{.Name}}={###}{{end -}}
	{{- if .Drop}}
	DROP TABLE IF EXISTS {{.Name}};
	{{end -}}
	CREATE TABLE IF NOT EXISTS {{.Name}} (
		{{range $i,$c:=.Rows -}}
		{{$c.Name}} {{$c.Type|dbType}} {{$c|defValue}} {{$c|isNull}} {{$c|seqTag}} comment '{{$c.Desc}}' {{if lt $i $count}},{{end}}
		{{end -}}{{.|indexs}}
	) ENGINE=InnoDB {{.|seqValue}} DEFAULT CHARSET=utf8mb4 COMMENT='{{.Desc}}'
  {{- if .PKG}}{###}{{end -}} `
