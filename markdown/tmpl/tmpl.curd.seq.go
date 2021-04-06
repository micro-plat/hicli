package tmpl

const MarkdownCurdSeqSql = `package sql

//SQLGetSEQ 获取序列
const SQLGetSEQ = {###}insert into sys_sequence_info (name,create_time) values (@name, now()){###}

//SQLClearSEQ 清除序列
const SQLClearSEQ = {###}delete from sys_sequence_info where seq_id < @seq_id{###}
`
