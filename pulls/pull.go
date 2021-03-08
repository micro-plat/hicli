package pulls

import (
	"fmt"

	"github.com/lib4dev/cli/cmds"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/gitlabs"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "pull",
			Usage: "拉取最新",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "branch,b",
					Usage: "分支名称",
				},
			},
			Action: pull,
		})
}

var pullBranch = []string{"master", "main"}

//pull 根据传入的路径(分组/仓库)拉取所有项目
func pull(c *cli.Context) (err error) {
	for _, b := range pullBranch {
		branch := types.GetString(c.String("branch"), b)
		reps, err := gitlabs.GetRepositories(c.Args().Get(0))
		if err != nil {
			continue
		}
		if len(reps) == 0 {
			return fmt.Errorf("没有需要拉取的项目")
		}
		for _, rep := range reps {
			if !rep.Exists() {
				logs.Log.Infof("get clone %s %s", rep.FullPath, rep.GetLocalPath())
				if err := rep.Clone(); err != nil {
					logs.Log.Error(err)
					continue
				}
			}
			if err := rep.Pull(branch); err != nil {
				logs.Log.Error(err)
				continue
			}
		}
		break
	}
	return nil

}
