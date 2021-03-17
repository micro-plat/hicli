package update

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/micro-plat/hicli/gitlabs"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "update",
			Usage:  "更新服务",
			Action: update,
		})
}

//pull 根据传入的路径(分组/仓库)拉取所有项目
func update(c *cli.Context) (err error) {
	rs, err := gitlabs.GetRepositories("github.com/micro-plat/hicli")
	if err != nil {
		return err
	}
	for _, r := range rs {
		if err := r.Update(); err != nil {
			return err
		}
	}
	return nil
}
