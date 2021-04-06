package tmpl

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

func (t *Table) Mysql2Column(data types.XMap) (c *Row) {
	return &Row{
		Name:   data.GetString("column_name"),
		Len:    -1,
		LenStr: "",
		Type:   strings.ToLower(data.GetString("column_type")),
		Def:    getDefault(data.GetString("column_default")),
		IsNull: data.GetString("is_nullable"),
		Con:    getMysqlCon(data.GetString("con")),
		Desc:   data.GetString("column_comment"),
	}
}

func (t *Table) Oracle2Column(data types.XMap) (c *Row) {
	con := fmt.Sprintf("%s,%s", data.GetString("constraint_type"), data.GetString("index_name"))
	return &Row{
		Name:   strings.ToLower(data.GetString("column_name")),
		Len:    data.GetInt("data_length", -1),
		LenStr: data.GetString("data_length"),
		Type:   strings.ToLower(data.GetString("data_type")),
		Def:    getDefault(data.GetString("data_default")),
		IsNull: data.GetString("nullable"),
		Con:    getOracleCon(con),
		Desc:   strings.TrimSpace(strings.ReplaceAll(data.GetString("column_comments"), "|", "&#124;")),
	}
}

func getMysqlCon(con string) string {
	if strings.Contains(con, "PRIMARY KEY") {
		con = "PK"
	}
	if strings.Contains(con, "UNIQUE") {
		con = strings.ReplaceAll(con, "UNIQUE", "UNQ")
	}
	return con
}

// C　　　　　　Check constraint on a table
// P　　　　　　Primary key
// U　　　　　　Unique key
// R　　　　　　Referential integrity
// V　　　　　　With check option, on a view
// O　　　　　　With read only, on a view
// H　　　　　　Hash expression
// F　　　　　　Constraint that involves a REF column
// S　　　　　　Supplemental loggin
func getOracleCon(s string) string {
	cons := strings.Split(s, ",")
	con := []string{}
	for _, v := range cons {
		if strings.HasPrefix(v, "C(") { //检查约束，跳过
			continue
		}
		if strings.HasPrefix(v, "P(") {
			con = append(con, "PK")
			continue
		}
		if strings.HasPrefix(v, "U(") {
			con = append(con, strings.Replace(strings.Replace(v, "U(", "UNQ(", 1), "|", ",", 1))
			continue
		}
		if strings.HasPrefix(v, "IDX(") {
			con = append(con, strings.Replace(v, "|", ",", 1))
		}
	}
	return strings.Join(con, ",")
}

func getDefault(str string) string {
	s := strings.TrimSpace(strings.ReplaceAll(str, "|", "&#124;"))
	s = strings.TrimPrefix(s, `'`)
	s = strings.TrimSuffix(s, `'`)
	return s
}
