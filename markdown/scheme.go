package markdown

import (
	"fmt"

	"github.com/micro-plat/hicli/markdown/const/enums"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/urfave/cli"
)

const gofile = "gofile"

//createScheme 生成数据库结构
func createScheme(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}
	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定输出路径")
	}

	filePath := c.Args().First()
	var tbs *tmpl.Tables
	tbs, err = tmpl.Markdowns2DB(filePath)
	if err != nil {
		return err
	}

	//设置包名称
	if c.Bool(gofile) {
		tbs.SetPkg(c.Args().Get(1))
	}
	//是否删除表
	tbs.DropTable(c.Bool("drop"))
	//过滤数据表
	tbs.FilterByKW(c.String("table"))
	tbs.Exclude()
	tbs.BuildSEQFile(c.Bool("seqfile"))

	//循环创建表
	outpath := c.Args().Get(1)
	for _, tb := range tbs.Tbs {
		//创建文件
		path := tmpl.GetSchemePath(outpath, tb.Name, c.Bool(gofile))

		//翻译文件
		content, err := tmpl.Translate(tmpl.SQLTmpl, dbtp, tb)
		if err != nil {
			return err
		}
		fs, err := tmpl.Create(path, c.Bool("cover"))
		if err != nil {
			return err
		}
		logs.Log.Info("生成文件:", path)
		if _, err := fs.Write([]byte(content)); err != nil {
			return err
		}
	}
	if tbs.SEQFile {
		content, err := tmpl.Translate(tmpl.MarkdownCurdSeqInstallSQL, dbtp, tbs)
		if err != nil {
			return err
		}
		path := tmpl.GetSEQFilePath(outpath, c.Bool(gofile))
		fs, err := tmpl.Create(path, c.Bool("cover"))
		if err != nil {
			return err
		}
		logs.Log.Info("生成文件:", path)
		fs.WriteString(content)
		fs.Close()
	}

	//生成安装文件
	if c.Bool(gofile) {
		content, err := tmpl.Translate(tmpl.InstallTmpl, dbtp, tbs)
		if err != nil {
			return err
		}
		path := tmpl.GetInstallPath(outpath)
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

//createDiff 生成数据库结构
func createDiff(c *cli.Context) (err error) {
	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定需要对比的源md文件和目标md文件")
	}

	sourcePath := c.Args().First()
	sourceTbs, err := tmpl.Markdowns2DB(sourcePath)
	if err != nil {
		return err
	}
	targetPath := c.Args().Get(1)
	targetTbs, err := tmpl.Markdowns2DB(targetPath)
	if err != nil {
		return err
	}

	//过滤表
	sourceTbs.FilterByKW(c.String("table"))
	sourceTbs.Exclude()

	outpath := c.Args().Get(2)

	diff := sourceTbs.Diff(targetTbs)

	if c.Bool(gofile) {
		diff.SetPkg(outpath)
	}

	for _, tb := range diff.Tbs {
		//创建文件
		path := tmpl.GetSchemePath(outpath, tb.Name, c.Bool(gofile))
		sqlTmpl := ""
		switch tb.Operation {
		case enums.DiffDelete:
			sqlTmpl = tmpl.DiffSQLDeleteTmpl
		case enums.DiffInsert:
			sqlTmpl = tmpl.DiffSQLInsertTmpl
		case enums.DiffModify:
			sqlTmpl = tmpl.DiffSQLModifyTmpl
		}
		//翻译文件
		content, err := tmpl.Translate(sqlTmpl, dbtp, tb)
		if err != nil {
			return err
		}
		fs, err := tmpl.Create(path, c.Bool("cover"))
		if err != nil {
			return err
		}
		logs.Log.Info("生成文件:", path)
		if _, err := fs.Write([]byte(content)); err != nil {
			return err
		}
	}

	//生成安装文件
	if c.Bool(gofile) {
		content, err := tmpl.Translate(tmpl.InstallTmpl, dbtp, diff)
		if err != nil {
			return err
		}
		path := tmpl.GetInstallPath(outpath)
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
