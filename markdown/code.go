package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"
)

func showEnitfy() func(c *cli.Context) (err error) {
	return showCode("entity")
}

func showField() func(c *cli.Context) (err error) {
	return showFiledCode("field")
}

//showCode 生成代码语句
func showCode(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//获取生成文件有关路径
		root := c.Args().Get(1)
		projectPath := utils.GetProjectPath(root)

		//读取文件
		tb, err := tmpl.Markdowns2DB(c.Args().First())
		if err != nil {
			return err
		}
		//过滤数据表
		tb.FilterByKW(c.String("table"))

		script := entityMap[tp]

		for _, tb := range tb.Tbs {
			path := tmpl.GetFileName(fmt.Sprintf("%s/modules/const/%s", projectPath, tp), tb.Name, "")
			tb.SetPkg(path)
			//翻译文件
			content, err := tmpl.Translate(script, dbtp, tb)
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
		}

		return nil
	}
}

func showFiledCode(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//获取field配置
		root := c.Args().Get(1)
		filedPath := tmpl.GetFieldConfPath(root)
		if filedPath == "" {
			return
		}
		projectPath := utils.GetProjectPath(root)
		basePath := utils.GetProjectBasePath(projectPath)

		//读取文件
		tbs, err := tmpl.Markdowns2DB(c.Args().First())
		if err != nil {
			return fmt.Errorf("处理markdown文件表格出错:%+v", err)
		}

		//过滤数据表
		tbs.FilterByKW(c.String("table"))

		for _, tb := range tbs.Tbs {
			//设置项目目录
			tb.SetBasePath(basePath)

			//保存的动态配置
			err := tmpl.NewFieldConf(tb).SaveConf(filedPath)
			if err != nil {
				logs.Log.Error(err)
			}
		}

		confs, err := tmpl.GetFieldConf(filedPath)
		if err != nil {
			return err
		}

		//翻译文件
		content, err := tmpl.Translate(entityMap[tp], dbtp, confs)
		if err != nil {
			return err
		}
		if !c.Bool("w2f") {
			logs.Log.Info(content)
			return nil
		}

		//生成文件
		path := tmpl.GetFileName(fmt.Sprintf("%s/modules/const/%s", projectPath, tp), "field", "")
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

var entityMap = map[string]string{
	"entity": tmpl.EntityTmpl,
	"field":  tmpl.FieldsTmpl,
}
