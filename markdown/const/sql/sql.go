package sql

// QueryMysql .
const QueryMysqlColumnInfo = `
select
c.table_name,
c.column_name,
c.column_default,
c.is_nullable,
c.column_type,
c.column_key,
c.column_comment,
c.ordinal_position,
t.table_comment,
GROUP_CONCAT(CONCAT(if(co.constraint_type is null,'IDX',co.constraint_type),'(',s.index_name,',',s.SEQ_IN_INDEX,')')) con
from
information_schema.columns c
left join information_schema.tables t on 
                (c.table_name = t.table_name and t.table_schema =c.table_schema)
left join information_schema.statistics s on 
                ( s.column_name = c.column_name and s.table_name = c.table_name and s.table_schema =c.table_schema  )
left join information_schema.key_column_usage kc on 
                ( kc.table_name = s.table_name and kc.column_name =s.column_name and kc.table_schema =s.table_schema and s.index_name=kc.constraint_name )
left join information_schema.table_constraints co on 
                ( kc.constraint_name = co.constraint_name and co.table_name = kc.table_name and co.table_schema =kc.table_schema ) 
where
c.table_schema =@schema
group by 
c.table_name,
c.column_name,
c.column_default,
c.is_nullable,
c.column_type,
c.column_key,
c.column_comment,
c.ordinal_position,
t.table_comment
order by c.table_name,c.ordinal_position
`

// GetAllTableNameInOracle .
const GetAllTableNameInOracle = `
select
  ub.table_name
from
  user_tables ub
order by 
  ub.table_name
`

// GetSingleTableInfoInOracle .
const GetSingleTableInfoInOracle = `SELECT a.table_name,
a.column_name,
a.data_type,
CASE upper(a.data_type)
  WHEN upper('number') THEN
   decode(a.data_scale,
		  0,
		  to_char(a.data_precision),
		  a.data_precision || ',' || a.data_scale)
  WHEN upper('date') THEN
   ''
  ELSE
   to_char(a.data_length)
END data_length,
a.data_precision,
a.data_scale,
a.nullable,
a.data_default data_default,
b.comments column_comments,
c.comments table_comments,
nvl(d.constraint_type,'') constraint_type,
nvl(e.index_name,'') index_name
FROM user_tab_columns a
LEFT JOIN user_col_comments b ON (a.table_name = b.table_name AND
							a.column_name = b.column_name)
LEFT JOIN user_tab_comments c ON a.table_name = c.table_name
LEFT JOIN (SELECT cu.table_name,
			 cu.column_name,
			 wm_concat(au.constraint_type || '(' ||
					   au.constraint_name || '|' ||
					   nvl(cu.position, 0) || ')') constraint_type
		FROM user_cons_columns cu
		LEFT JOIN user_constraints au ON (cu.constraint_name =
										 au.constraint_name)
	   WHERE cu.table_name = @table_name
	   GROUP BY cu.table_name, cu.column_name) d ON d.column_name = a.column_name
LEFT JOIN (select t.column_name,
			 wm_concat('IDX(' || t.index_name || '|' ||
					   nvl(t.column_position, 0) || ')') index_name
		from user_ind_columns t, user_indexes i
	   where t.table_name = @table_name
		 and t.index_name = i.index_name
		 and i.uniqueness = 'NONUNIQUE'
	   group by t.column_name) e on e.column_name = a.column_name
WHERE a.table_name = @table_name
ORDER BY a.column_id
`

const GetMysqlColumnInfo = `
select
c.table_name,
GROUP_CONCAT(c.column_name) column_name
from
information_schema.columns c
where
c.table_schema = @schema
group by 
c.table_name
order by c.table_name,c.ordinal_position
`

const ExportMysqlData = `SELECT #column_name value from #table_name t`
