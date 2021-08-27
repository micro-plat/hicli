package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"
)

func createCurd() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		//创建sql
		if err = showSQL("curd")(c); err != nil {
			return err
		}
		//创建数据库引擎引入文件
		if err = createConstFile("driver")(c); err != nil {
			return err
		}
		if isOracle {
			return
		}
		//创建序列sql
		if err = createConstFile("seq")(c); err != nil {
			return err
		}
		//创建序列表安装文件
		if err = createConstFile("seq.install.go")(c); err != nil {
			return err
		}
		///创建序列表安装sql
		if err = createConstFile("seq.install.sql")(c); err != nil {
			return err
		}
		return nil
	}
}

func showSelect() func(c *cli.Context) (err error) {
	return showSQL("select")
}

func showUpdate() func(c *cli.Context) (err error) {
	return showSQL("update")
}

func showInsert() func(c *cli.Context) (err error) {
	return showSQL("insert")
}

//showSQL 生成SQL语句
func showSQL(sqlType string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//获取相关路径
		root := c.Args().Get(1)
		projectPath := utils.GetProjectPath(root)

		//读取文件
		tb, err := tmpl.Markdowns2DB(c.Args().First())
		if err != nil {
			return err
		}

		allTables := tb.Tbs
		for _, tb := range tb.Tbs {
			tb.SetAllTables(allTables)
			tb.DisposeELTab()
			tb.DispostELBtnList()
		}

		//过滤数据表
		tb.FilterByKW(c.String("table"))

		for _, tb := range tb.Tbs {
			path := tmpl.GetFileName(fmt.Sprintf("%s/modules/const/sql", projectPath), tb.Name, fmt.Sprintf("%s.", dbtp))
			//根据关键字过滤
			tb.FilterRowByKW(c.String("kw"))
			tb.DBType = dbtp
			tb.SetPkg(path)

			//翻译文件
			content, err := tmpl.Translate(sqlMap[sqlType], dbtp, tb)
			if err != nil {
				return fmt.Errorf("翻译%s模板出错:%+v", sqlType, err)
			}
			if !c.Bool("w2f") {
				logs.Log.Info(content)
				return nil
			}
			//生成文件
			fs, err := tmpl.Create(path, c.Bool("cover"))
			if err != nil {
				return err
			}
			logs.Log.Info("生成文件:", path)
			fs.WriteString(content)
			fs.Close()
		}
		return nil
	}
}

func createConstFile(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//获取相关路径
		root := c.Args().Get(1)
		projectPath := utils.GetProjectPath(root)

		path := tmpl.GetFileName(fmt.Sprintf("%s/modules/const/sql", projectPath), sqlPathMap[tp], dbtp)

		//文件存在则不生成
		if tmpl.PathExists(path) {
			return
		}
		//翻译文件
		content, err := tmpl.Translate(sqlMap[tp], dbtp, map[string]interface{}{
			"BasePath": utils.GetProjectBasePath(projectPath),
		})
		if err != nil {
			return err
		}
		if !c.Bool("w2f") {
			logs.Log.Info(content)
			return nil
		}
		//生成文件
		fs, err := tmpl.Create(path, c.Bool("cover"))
		if err != nil {
			return err
		}
		logs.Log.Info("生成文件:", path)
		fs.WriteString(content)
		fs.Close()
		return nil
	}
}

var sqlMap = map[string]string{
	"insert":          tmpl.InsertSingle,
	"update":          tmpl.UpdateSingle,
	"select":          tmpl.SelectSingle,
	"curd":            tmpl.MarkdownCurdSql,
	"driver":          tmpl.MarkdownCurdDriverSql,
	"seq":             tmpl.MarkdownCurdSeqSql,
	"seq.install.go":  tmpl.MarkdownCurdSeqInstallGO,
	"seq.install.sql": tmpl.MarkdownCurdSeqInstallSQL,
}

var sqlPathMap = map[string]string{
	"driver":          "",
	"seq":             ".seq.info",
	"seq.install.go":  "/install",
	"seq.install.sql": "_/sys_sequence_info",
}
