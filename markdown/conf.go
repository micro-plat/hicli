package markdown

import (
	"fmt"
	"path"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/app"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/hicli/markdown/ui"
	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"
)

//createVueRouter 创建vue路由
func createVueRouter() func(c *cli.Context) (err error) {
	return createConf("vue.router")
}

//createVueMenus 创建vue菜单
func createVueMenus() func(c *cli.Context) (err error) {
	return createConf("vue.menus")
}

func createGORouter() func(c *cli.Context) (err error) {
	return createGo("init.go")
}

func createConf(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//获取配置
		root := c.Args().Get(1)
		confPath := tmpl.GetWebConfPath(root)
		if confPath == "" {
			return
		}
		confs, err := tmpl.GetSnippetConf(confPath)
		if err != nil {
			return err
		}
		//翻译
		content, err := tmpl.Translate(confMap[tp], "", confs)
		if err != nil {
			return err
		}

		//获取有关路径
		projectPath := utils.GetProjectPath(root)
		webPath, webSrcPath := utils.GetWebSrcPath(projectPath)
		if webSrcPath == "" {
			return
		}
		path := path.Join(webPath, confPathMap[tp])
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

func createGo(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//获取配置
		root := c.Args().Get(1)
		confPath := tmpl.GetGoConfPath(root)
		if confPath == "" {
			return
		}
		confs, err := tmpl.GetSnippetConf(confPath)
		if err != nil {
			return err
		}

		projectPath := utils.GetProjectPath(root)
		//翻译内容
		content, err := tmpl.Translate(confMap[tp], "", map[string]interface{}{
			"GOMOD":       utils.GetGOMOD(),
			"ProjectPath": projectPath,
			"Confs":       confs,
		})
		if err != nil {
			return err
		}

		//生成文件
		path := path.Join(projectPath, "init.go")
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

var confMap = map[string]string{
	"vue.router": ui.SnippetSrcRouterIndexJS,
	"vue.menus":  ui.SrcMenusJson,
	"init.go":    app.SnippetTmplConfGo,
}

var confPathMap = map[string]string{
	"vue.router": "src/router/index.js",
	"vue.menus":  "public/menus.json",
}
