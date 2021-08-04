// +build !oracle

package tmpl

const SQLTmpl = `
{{- if .PKG}}package {{.PKG}}
{{end -}}

{{$count:=.Rows|maxIndex -}}

{{- if .PKG}} 
//{{.Name}} {{.Desc}}
const {{.Name}} = {###}{{end -}}
	{{- if .Drop}}
	DROP TABLE IF EXISTS {{.Name}};
	{{end -}}
	CREATE TABLE IF NOT EXISTS {{.Name}} (
		{{range $i,$c:=.Rows -}}
		{{$c.Name}} {{$c.Type|dbType}} {{$c|defValue}} {{$c|isNull}} {{$c|seqTag}} comment '{{$c.Desc}}' {{if lt $i $count}},{{end}}
		{{end -}}{{.|indexs}}
	) ENGINE=InnoDB {{.|seqValue}} DEFAULT CHARSET=utf8mb4 COMMENT='{{.Desc}}'
  {{- if .PKG}}{###}{{end -}};`

const DiffSQLInsertTmpl = `
{{- if .PKG}}package {{.PKG}}
{{end -}}

{{$count:=.Rows|maxIndex -}}

{{- if .PKG}} 
//{{.Name}} {{.Desc}}
const {{.Name}} = {###}{{end -}}
	CREATE TABLE IF NOT EXISTS {{.Name}} (
		{{range $i,$c:=.Rows -}}
		{{$c.Name}} {{$c.Type|dbType}} {{$c|defValue}} {{$c|isNull}} {{$c|seqTag}} comment '{{$c.Desc}}' {{if lt $i $count}},{{end}}
		{{end -}}{{.|indexs}}
	) ENGINE=InnoDB {{.|seqValue}} DEFAULT CHARSET=utf8mb4 COMMENT='{{.Desc}}'
	{{- if .PKG}}{###}{{end -}};`

const DiffSQLDeleteTmpl = `
{{- if .PKG}}package {{.PKG}}
{{end -}}
{{- if .PKG}} 
//{{.Name}} {{.Desc}}
const {{.Name}} = {###}{{end -}}
	DROP TABLE IF EXISTS {{.Name}};{{- if .PKG}}{###}{{end -}};`

const DiffSQLModifyTmpl = `
{{- if .PKG}}package {{.PKG}}
{{end -}}
{{- if .PKG}} 
//{{.Name}} {{.Desc}}
const {{.Name}} = {###}{{end -}}
{{- range $i,$c:=.DiffRows}}
{{- if (eq $c.Operation -1)}}
-- 删除字段 {{$c.Name}} 
ALTER TABLE {{$.Name}} drop COLUMN {{$c.Name}};
{{- else if (eq $c.Operation 1)}}
-- 新增字段 {{$c.Name}} 
ALTER TABLE {{$.Name}} add COLUMN {{$c.Name}} {{$c.Type|dbType}} {{$c|defValue}} {{$c|isNull}} {{$c|seqTag}} comment '{{$c.Desc}}';
{{- else if (eq $c.Operation 2)}}
-- 修改字段 {{$c.Name}} 
ALTER TABLE {{$.Name}} MODIFY {{$c.Name}} {{$c.Type|dbType}} {{$c|defValue}} {{$c|isNull}} {{$c|seqTag}} comment '{{$c.Desc}}';
{{- end}}
{{- end}}
{{- range $i,$c:=.DiffIndexs}}
{{- if and (eq $c.Operation -1) ($c|isPK)}}
-- 删除主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP PRIMARY KEY;
{{- else if and (eq $c.Operation 1) ($c|isPK)}}
-- 新增主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} ADD CONSTRAINT symbol {{$c|indexStr}};
{{- else if and (eq $c.Operation 2) ($c|isPK)}}
-- 修改主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP PRIMARY KEY;
ALTER TABLE {{$.Name}} ADD CONSTRAINT symbol {{$c|indexStr}};
{{- else if and (eq $c.Operation -1) (or ($c|isIndex) ($c|isUNQ))}}
-- 删除索引 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP INDEX {{$c.Name}};
{{- else if and (eq $c.Operation 1) (or ($c|isIndex) ($c|isUNQ))}}
-- 新增索引 {{$c.Name}} 
ALTER TABLE {{$.Name}} ADD CONSTRAINT symbol {{$c|indexStr}};
{{- else if and (eq $c.Operation 2) (or ($c|isIndex) ($c|isUNQ))}}
-- 修改索引 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP INDEX {{$c.Name}};
ALTER TABLE {{$.Name}} ADD CONSTRAINT symbol {{$c|indexStr}};
{{- end}}
{{- end}}
{{- if .PKG}}{###}{{end -}};`
