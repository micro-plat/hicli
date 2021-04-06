package clones

import (
	"github.com/lib4dev/cli/cmds"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/gitlabs"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "clone",
			Usage:  "克隆项目",
			Action: clone,
		})
}

//clone 根据传入的路径(分组/仓库)拉取所有仓库
func clone(c *cli.Context) (err error) {
	reps, err := gitlabs.GetRepositories(c.Args().Get(0))
	if err != nil {
		return err
	}
	for _, rep := range reps {
		if err := rep.Clone(); err != nil {
			logs.Log.Error(err)
		}
	}
	return nil

}
