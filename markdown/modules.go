package markdown

import (
	"fmt"
	"path"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"
)

func createModulesSeq() func(c *cli.Context) (err error) {
	return createModules("seq")
}

func createModules(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		root := c.Args().Get(1)
		projectPath := utils.GetProjectPath(root)
		basePath := utils.GetProjectBasePath(projectPath)

		//判断文件是否存在
		path := path.Join(projectPath, "modules/db/mysql.seq.info.go")
		if tmpl.PathExists(path) {
			return
		}

		//翻译
		content, err := tmpl.Translate(modulesMap[tp], "", map[string]interface{}{
			"BasePath":    basePath,
			"ProjectPath": projectPath,
		})
		if err != nil {
			return err
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

var modulesMap = map[string]string{
	"seq": tmpl.ModulesDBSeqTmpl,
}
