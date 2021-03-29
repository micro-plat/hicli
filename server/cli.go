package server

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "server",
			Usage:  "运行后端应用程序",
			Action: runServer(),
			Flags: []cli.Flag{
				cli.StringFlag{Name: "run", Required: false, Usage: `-应用程序启动参数`},
				cli.StringFlag{Name: "install", Required: false, Usage: `-go install参数`},
				cli.BoolFlag{Name: "work,w", Required: false, Usage: `-以当前路径为工作目录`},
			},
		})
}
