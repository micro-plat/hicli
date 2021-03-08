package resets

import (
	"github.com/lib4dev/cli/cmds"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/gitlabs"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "reset",
			Usage: "重置代码",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "branch,b",
					Usage: "分支名称",
				},
			},
			Action: reset,
		})
}

//reset 根据传入的路径(分组/仓库)重置所有本地修改
func reset(c *cli.Context) (err error) {
	reps, err := gitlabs.GetRepositories(c.Args().Get(0))
	if err != nil {
		return err
	}
	for _, rep := range reps {
		branch := types.GetString(c.String("branch"), "master")
		if err := rep.Reset(branch); err != nil {
			logs.Log.Error(err)
		}

	}
	return nil

}
