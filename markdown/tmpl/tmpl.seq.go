// +build !oracle

package tmpl

const CreateSEQTable = `
{{- if .PKG}}package {{.PKG}}
{{end -}}

{{- if .PKG}} 
//SEQ_IDS 序列生成表
const SEQ_IDS={###}{{end -}}
{{- if .Drop}}
DROP TABLE IF EXISTS seq_ids;
{{end -}}
CREATE TABLE  seq_ids (
id bigint  not null auto_increment comment '编号' ,
name varchar(64)  not null  comment '名称' ,
create_time datetime default current_timestamp not null  comment '创建时间' ,	
primary key (id)
) ENGINE=InnoDB auto_increment = 100 DEFAULT CHARSET=utf8 COMMENT='序列信息表'
{{- if .PKG}}{###}{{end -}} `
