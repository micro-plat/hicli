package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/app"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"
)

func createApp(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定项目名称")
	}
	//创建项目
	err = app.CreateApp(c.Args().First(), c.Bool("sso"))
	return
}

func createServiceBlock() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		//生成services
		if err := createBlockCode("service")(c); err != nil {
			return err
		}
		//生成sql
		if err := createCurd()(c); err != nil {
			return err
		}
		//生成field.go
		if err := showField()(c); err != nil {
			return err
		}
		//生成conf.go
		if err := createGORouter()(c); err != nil {
			return err
		}
		if c.Bool("oracle") {
			return
		}
		//生成modules/db/seq
		return createModulesSeq()(c)
	}
}

func createBlockCode(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		dbtp := tmpl.MYSQL
		if c.Bool("oracle") {
			dbtp = tmpl.ORACLE
		}

		//获取生成文件有关路径
		root := c.Args().Get(1)
		confPath := tmpl.GetGoConfPath(root)
		fieldPath := tmpl.GetFieldConfPath(root)
		projectPath := utils.GetProjectPath(root)
		basePath := utils.GetProjectBasePath(projectPath)

		//读取文件
		tbs, err := tmpl.Markdown2DB(c.Args().First())
		if err != nil {
			return fmt.Errorf("处理markdown文件表格出错:%+v", err)
		}

		allTables := tbs.Tbs
		for _, tb := range tbs.Tbs {
			tb.SetAllTables(allTables)
			tb.DisposeTabTables()
		}

		//过滤数据表
		tbs.FilterByKW(c.String("table"))

		for _, tb := range tbs.Tbs {
			//设置项目目录
			tb.SetBasePath(basePath)

			//保存的动态配置
			err := tmpl.NewSnippetConf(tb).SaveConf(confPath)
			if err != nil {
				logs.Log.Error(err)
			}
			err = tmpl.NewFieldConf(tb).SaveConf(fieldPath)
			if err != nil {
				logs.Log.Error(err)
			}

			//根据关键字过滤
			path := tmpl.GetFilePath(fmt.Sprintf("%s/services", projectPath), tb.Name, "go")
			tb.FilterRowByKW(c.String("kw"))
			tb.SetPkg(path)

			//翻译文件
			content, err := tmpl.Translate(appCodeMap[tp], dbtp, tb)
			if err != nil {
				return fmt.Errorf("翻译%s模板出错:%+v", tp, err)
			}

			//终端输出内容
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

func createEnums() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if err := createEnum()(c); err != nil {
			return err
		}
		return createGORouter()(c)
	}
}

func createEnum() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//获取生成文件有关路径
		root := c.Args().Get(1)
		projectPath := utils.GetProjectPath(root)

		//读取文件
		tbs, err := tmpl.Markdown2DB(c.Args().First())
		if err != nil {
			return fmt.Errorf("处理markdown文件表格出错:%+v", err)
		}
		//过滤数据表
		tbs.FilterByKW(c.String("table"))

		//设置包名
		path := tmpl.GetFilePath(fmt.Sprintf("%s/services/system", projectPath), "system.enums", "go")
		tbs.SetPkg(path)

		tml := app.TmplEnumsHandler
		if c.Bool("dds") {
			tml = app.TmplEnumsHandlerDDS
		}

		//翻译文件
		dbtp := tmpl.MYSQL
		if c.Bool("oracle") {
			dbtp = tmpl.ORACLE
		}
		content, err := tmpl.Translate(tml, dbtp, tbs)
		if err != nil {
			return fmt.Errorf("翻译%s模板出错:%+v", "enums", err)
		}
		if !c.Bool("w2f") {
			logs.Log.Info(content)
			return nil
		}

		//保存配置内容
		basePath := utils.GetProjectBasePath(projectPath)
		confPath := tmpl.GetGoConfPath(root)
		tb := &tmpl.Table{
			Name:     "_system_enums",
			BasePath: basePath,
		}
		err = tmpl.NewSnippetConf(tb).SaveConf(confPath)
		if err != nil {
			logs.Log.Error(err)
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

var appCodeMap = map[string]string{
	"service": app.TmplServiceHandler,
}
