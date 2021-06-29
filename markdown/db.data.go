package markdown

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/const/sql"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

//exportDBData 导出数据库数据
func exportDBData(c *cli.Context) (err error) {

	str := strings.Split(c.String("dbstr"), ":")
	provider := str[0]
	if _, ok := dbDataProviderFunc[provider]; !ok {
		return fmt.Errorf("连接串有误%s", c.String("dbstr"))
	}

	//获取对应表
	f := dbDataProviderFunc[provider]
	connstr := strings.TrimPrefix(c.String("dbstr"), fmt.Sprintf("%s:", provider))
	values, err := f(provider, connstr, c.String("table"))
	if err != nil {
		return err
	}
	list := ""
	for _, v := range values {
		if v == "" {
			continue
		}
		list += "	" + fmt.Sprintf(`"%s",`, v) + "\n"
	}

	//循环创建表
	content := ""
	//翻译文件
	ct, err := tmpl.Translate(tmpl.MdDBDataTPL, provider, map[string]interface{}{
		"PKG":      provider,
		"DataList": strings.TrimRight(list, "\n"),
	})
	if err != nil {
		return err
	}
	content += ct

	if !c.Bool("w2f") {
		logs.Log.Info(content)
		return nil
	}

	//生成文件
	tempArr := strings.Split(connstr, "/")
	tableScheme := tempArr[len(tempArr)-1]
	if provider == "oracle" {
		tableScheme = tempArr[0]
	}

	path := filepath.Join(c.Args().Get(0), fmt.Sprintf("./%s.data.go", tableScheme))
	fs, err := tmpl.Create(path, c.Bool("cover"))
	if err != nil {
		return err
	}
	logs.Log.Info("生成文件:", path)
	fs.WriteString(content)
	fs.Close()

	return nil
}

var dbDataProviderFunc = map[string]func(string, string, string) ([]string, error){
	"mysql": exportMysqlData,
}

//exportMysqlData 导出mysql数据 mysql:root:xxxx@tcp(192.168.0.36:3306)/sms_test
func exportMysqlData(provider, connstr, subTable string) (values []string, err error) {
	tempArr := strings.Split(connstr, "/")
	tableScheme := tempArr[len(tempArr)-1]
	dbObj, err := db.NewDB(provider, connstr, 20, 10, 20000)
	if err != nil {
		return nil, fmt.Errorf("NewDB:%+v,provider:%s,connstr:%s", err, provider, connstr)
	}
	tbs, err := dbObj.Query(sql.GetMysqlColumnInfo, map[string]interface{}{
		"schema": tableScheme,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMysqlColumnInfo:%+v,provider:%s,connstr:%s", err, provider, connstr)
	}

	values = make([]string, 0)

	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, v := range tbs {
		if !strings.Contains(v.GetString("table_name"), subTable) {
			continue
		}
		wg.Add(1)
		go func(dbObj db.IDB, v types.XMap) {
			defer wg.Done()
			r := querySQL(dbObj, v)
			lock.Lock()
			defer lock.Unlock()
			values = append(values, r)
		}(dbObj, v)
	}
	wg.Wait()

	return
}

func querySQL(db db.IDB, table types.XMap) string {

	columns := strings.Split(table.GetString("column_name"), ",")
	if len(columns) < 1 {
		return ""
	}
	s := ""
	for _, v := range columns {
		s += fmt.Sprintf(`case when %s is null then "null," else concat("'",%s,"',") end,`, v, v)
	}
	queryStr := fmt.Sprintf(`concat(%s)`, strings.TrimRight(s, ","))
	rows, err := db.Query(sql.ExportMysqlData, map[string]interface{}{
		"table_name":  table.GetString("table_name"),
		"column_name": queryStr})
	if err != nil {
		panic(err)
	}
	if rows.IsEmpty() {
		return ""
	}
	r := fmt.Sprintf("INSERT INTO %s(%s) VALUES ", table.GetString("table_name"), table.GetString("column_name"))
	for _, v := range rows {
		r += fmt.Sprintf("(%s),", strings.TrimRight(v.GetString("value"), ","))
	}

	return strings.TrimRight(r, ",")
}
