package server

import (
	"fmt"
	"path/filepath"

	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"
)

func runServer() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		//判断项目是否存在
		projectPath := utils.GetProjectPath(c.Args().Get(0))
		if !utils.PathExists(filepath.Join(projectPath, "main.go")) {
			return fmt.Errorf("未指定的运行应用程序的项目路径:%s", projectPath)
		}

		//构建服务
		s, err := newServer(c, projectPath)
		if err != nil {
			return err
		}

		//服务启动
		go s.resume()

		//服务关闭
		err = s.close()

		return
	}
}
