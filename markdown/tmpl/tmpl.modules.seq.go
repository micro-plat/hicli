package tmpl

const ModulesDBSeqTmpl = `
package db

import (
	"fmt"
	"{{.BasePath}}/modules/const/sql"

	"github.com/micro-plat/lib4go/db"
)

// GetNewID 获取新ID
func GetNewID(db db.IDBExecuter, SQLGetSEQ string, imap map[string]interface{}) (taskID int64, err error) {
	id, row, err := db.Executes(SQLGetSEQ, imap)
	if err != nil || row < 1 {
		return 0, fmt.Errorf("获取批次编号失败 %v", err)
	}

	if id%1000 == 100 {
		_, err = db.Execute(sql.SQLClearSEQ, map[string]interface{}{"seq_id": id})
		if err != nil {
			return 0, fmt.Errorf("清理序列数据失败%v", err)
		}
	}
	return id, nil
}
`
