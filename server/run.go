package server

import (
	"fmt"
	"path/filepath"

	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"
)

func runServer() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		//path
		pPath, wPath, err := checkPath(c)
		if err != nil {
			return err
		}

		//构建服务
		s, err := newServer(c, pPath, wPath)
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

func checkPath(c *cli.Context) (projectPath, watchPath string, err error) {

	//判断项目是否存在
	projectPath = utils.GetProjectPath(c.Args().Get(0))
	if !utils.PathExists(filepath.Join(projectPath, "main.go")) {
		return "", "", fmt.Errorf("未指定的运行应用程序的项目路径:%s", projectPath)
	}

	//指定监控当前目录
	watchPath = projectPath
	if c.Bool("work") {
		watchPath = utils.GetProjectPath("./")
	}

	return
}
