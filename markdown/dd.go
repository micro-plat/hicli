package markdown

import (
	"fmt"
	"path/filepath"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/const/sql"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/lib4go/db"
	"github.com/urfave/cli"
)

//createScheme 生成数据库结构
func createDataDic(c *cli.Context) (err error) {

	str := strings.Split(c.String("dbstr"), ":")
	provider := str[0]
	if _, ok := dbProviderFunc[provider]; !ok {
		return fmt.Errorf("连接串有误%s", c.String("dbstr"))
	}

	//获取对应表
	f := dbProviderFunc[provider]
	connstr := strings.TrimPrefix(c.String("dbstr"), fmt.Sprintf("%s:", provider))
	tbs, err := f(provider, connstr)
	if err != nil {
		return err
	}

	//循环创建表
	content := ""
	for _, tb := range tbs {
		//翻译文件
		ct, err := tmpl.Translate(tmpl.MdDictionaryTPL, dbtp, tb)
		if err != nil {
			return err
		}
		content += ct
	}

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

	path := filepath.Join(c.Args().Get(0), fmt.Sprintf("./%s.%s.md", tableScheme, provider))
	fs, err := tmpl.Create(path, c.Bool("cover"))
	if err != nil {
		return err
	}
	logs.Log.Info("生成文件:", path)
	fs.WriteString(content)
	fs.Close()

	return nil
}

var dbProviderFunc = map[string]func(string, string) ([]*tmpl.Table, error){
	"mysql":  generateMysqlMD,
	"oracle": generateOracleMD,
}

//generateMysqlMD 生成mysql的markdown文件 mysql:root:xxxx@tcp(192.168.0.36:3306)/sms_test
func generateMysqlMD(provider, connstr string) (tabs []*tmpl.Table, err error) {
	tempArr := strings.Split(connstr, "/")
	tableScheme := tempArr[len(tempArr)-1]
	dbObj, err := db.NewDB(provider, connstr, 20, 10, 20000)
	if err != nil {
		return nil, fmt.Errorf("createMysqlMD，NewDB出错，err:%+v,provider:%s,connstr:%s", err, provider, connstr)
	}
	datas, err := dbObj.Query(sql.QueryMysqlColumnInfo, map[string]interface{}{
		"schema": tableScheme,
	})
	if err != nil {
		return nil, fmt.Errorf("mysql(err:%v)", err)
	}
	if len(datas) < 1 {
		return nil, fmt.Errorf("mysql(未查询到相关信息)，scheme:%s，data:%v", tableScheme, datas)
	}

	d := map[string]*tmpl.Table{}
	for _, v := range datas {
		tableName := v.GetString("table_name")
		tableComment := v.GetString("table_comment")
		if _, ok := d[tableName]; !ok {
			d[tableName] = tmpl.NewTable(tableName, tableComment, "")
		}
		d[tableName].AddRow(d[tableName].Mysql2Column(v))
	}

	tbs := []*tmpl.Table{}
	for _, v := range d {
		tbs = append(tbs, v)
	}
	return tbs, nil
}

//generateMysqlMD 生成oracle的markdown文件 //oracle:sso/123456@orcl136
func generateOracleMD(provider, connstr string) (tabs []*tmpl.Table, err error) {
	dbObj, err := db.NewDB(provider, connstr, 20, 10, 20000)
	if err != nil {
		return nil, fmt.Errorf("createOracleMD，NewDB出错，err:%+v,provider:%s,connstr:%s", err, provider, connstr)
	}
	tableNames, err := dbObj.Query(sql.GetAllTableNameInOracle, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("oracle(err:%v)", err)
	}

	d := map[string]*tmpl.Table{}
	for _, v := range tableNames {
		tableName := v.GetString("table_name")
		datas, err := dbObj.Query(sql.GetSingleTableInfoInOracle, map[string]interface{}{
			"table_name": strings.ToUpper(tableName),
		})
		if err != nil {
			return nil, fmt.Errorf("oracle(err:%v)", err)
		}
		if _, ok := d[tableName]; !ok {
			d[tableName] = tmpl.NewTable(strings.ToLower(tableName), datas.Get(0).GetString("table_comments"), "")
		}
		for _, v2 := range datas {
			d[tableName].AddRow(d[tableName].Oracle2Column(v2))
		}
	}
	tbs := []*tmpl.Table{}
	for _, v := range d {
		tbs = append(tbs, v)
	}
	return tbs, nil
}
