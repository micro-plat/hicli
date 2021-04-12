package server

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "server",
			Usage: "运行应用程序",
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "run",
					Usage:  "运行并监控应用程序",
					Action: runServer(),
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "fixed,fd", Usage: `-指定服务名称与程序名称相同`},
						cli.StringFlag{Name: "registry,r", EnvVar: "registry", Usage: `-注册中心地址。格式：proto://host。如：zk://ip1,ip2  或 fs://../`},
						cli.StringFlag{Name: "plat,p", Usage: "-平台名称"},
						cli.StringFlag{Name: "system,s", Usage: "-系统名称,默认为当前应用程序名称"},
						cli.StringFlag{Name: "server-types,S", Usage: "-服务类型，有api,web,rpc,cron,mqc,ws。多个以“-”分割"},
						cli.StringFlag{Name: "cluster,c", Usage: "-集群名称，默认值为：prod"},
						cli.BoolFlag{Name: "debug,d", Usage: `-调试模式，打印更详细的系统运行日志，避免将详细的错误信息返回给调用方`},
						cli.StringFlag{Name: "trace,t", Usage: `-性能分析。支持:cpu,mem,block,mutex,web`},
						cli.StringFlag{Name: "tport,tp", Usage: `-性能分析服务端口号。用于trace为web模式时的端口号。默认：19999`},
						cli.StringFlag{Name: "mask,msk", Usage: `-子网掩码。多个网卡情况下根据mask获取本机IP`},
						cli.StringFlag{Name: "tags", Usage: "-go 安装和打包编译的tags"},
						cli.StringFlag{Name: "mod", Usage: "-go 安装和打包编译的mod"},
						cli.StringFlag{Name: "run", Required: false, Usage: `-应用程序启动参数`},
						cli.StringFlag{Name: "install", Required: false, Usage: `-go install参数`},
						cli.BoolFlag{Name: "work,w", Required: false, Usage: `-以当前路径为工作目录`},
					},
				},
			},
		})
}
