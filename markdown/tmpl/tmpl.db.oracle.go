// +build oracle

package tmpl

const SQLTmpl = `
{{- if .PKG}}package {{.PKG}}
{{end -}}

{{$count:=.Rows|maxIndex -}}

{{- if .PKG}} 
//{{.Name}} {{.Desc}}
const {{.Name}}={###}{{end -}}
	{{- if .Drop}}
	drop table {{.Name}};
	{{- end}}
	create table {{.Name}}(
		{{- range $i,$c:=.Rows}}
		{{$c.Name}} {{$c.Type|dbType}} {{$c|defValue}} {{$c|isNull}} {{if lt $i $count}},{{end}}
		{{- end}}
	);

	comment on table {{.Name}} is '{{.Desc}}';
	{{- range $i,$c:=.Rows}}
	comment on column {{$.Name}}.{{$c.Name}} is '{{$c.Desc}}';	
	{{- end}}

	{{.|indexs}}

	{{- range $i,$c:=.|seqs}}
	create sequence {{$c.seqname}}
	increment by {{$c.increment}}
	minvalue {{$c.min}}
	maxvalue {{$c.max}}
	start with {{$c.min}}
	cache 20;
	{{- end}}

	{{if .PKG}}{###}{{end}} 
	`
