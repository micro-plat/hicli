package markdown

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "app",
			Usage: "后端应用程序",
			Subcommands: cli.Commands{
				{
					Name:   "create",
					Usage:  "创建app应用",
					Action: createApp,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "sso", Usage: `-生成sso相关内容`},
					},
				},
				{
					Name:   "enums",
					Usage:  "创建enums",
					Action: createEnums(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
						cli.BoolFlag{Name: "dds", Usage: `-引用dds`},
					},
				},
				{
					Name:   "service",
					Usage:  "创建服务",
					Action: createServiceBlock(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
			},
		},
		cli.Command{
			Name:  "ui",
			Usage: "创建vue前端项目",
			Subcommands: cli.Commands{
				{
					Name:   "create",
					Usage:  "创建项目",
					Action: createUI,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "sso", Usage: `-生成sso相关内容`},
					},
				},
				{
					Name:   "clear",
					Usage:  "清理项目",
					Action: clear,
				},
				{
					Name:   "page",
					Usage:  "创建项目页面",
					Action: createPage,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
				{
					Name:   "list",
					Usage:  "生成列表代码",
					Action: createList(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
				{
					Name:   "detail",
					Usage:  "生成预览代码",
					Action: createDetail(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
				{
					Name:   "edit",
					Usage:  "生成预览代码",
					Action: createEdit(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
				{
					Name:   "add",
					Usage:  "生成预览代码",
					Action: createAdd(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
			},
		},
		cli.Command{
			Name:  "code",
			Usage: "实体类文件",
			Subcommands: cli.Commands{
				{
					Name:   "entity",
					Usage:  "创建实体类",
					Action: showEnitfy(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
					},
				},
				{
					Name:   "field",
					Usage:  "创建表字段列表",
					Action: showField(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
			},
		},
		cli.Command{
			Name:  "db",
			Usage: "数据库结构文件",
			Subcommands: cli.Commands{
				{
					Name:   "create",
					Usage:  "创建数据库结构文件",
					Action: createScheme,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "gofile,g", Usage: `-生成到gofile中`},
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.BoolFlag{Name: "drop,d", Usage: `-包含表删除语句`},
						cli.BoolFlag{Name: "seqfile,s", Usage: `-包含序列文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
				{
					Name:   "diff",
					Usage:  "创建数据库结构差异文件",
					Action: createDiff,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.BoolFlag{Name: "gofile,g", Usage: `-生成到gofile中`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
				{
					Name:   "data",
					Usage:  "导出数据库数据",
					Action: exportDBData,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbstr,db", Required: true, Usage: `-数据库连接串，参考hydra的db配置时的连接串`},
						cli.BoolFlag{Name: "gofile,g", Usage: `-生成到gofile中`},
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
				{
					Name:   "select",
					Usage:  "创建select语句",
					Action: showSelect(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
					},
				},
				{
					Name:   "update",
					Usage:  "创建update语句",
					Action: showUpdate(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
					},
				},
				{
					Name:   "insert",
					Usage:  "创建insert语句",
					Action: showInsert(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
					},
				},
				{
					Name:   "crud",
					Usage:  "创建crud语句",
					Action: createCurd(),
					Flags: []cli.Flag{
						cli.StringFlag{Name: "kw,k", Usage: `-约束字段`},
						cli.StringFlag{Name: "table,t", Usage: `-表名称`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
					},
				},
			},
		},
		cli.Command{
			Name:  "dic",
			Usage: "生成数据库对应数据字典",
			Subcommands: cli.Commands{
				{
					Name:   "create",
					Usage:  "创建数据字典",
					Action: createDataDic,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "dbstr,db", Required: true, Usage: `-数据库连接串，参考hydra的db配置时的连接串`},
						cli.BoolFlag{Name: "cover,v", Usage: `-文件已存在时自动覆盖`},
						cli.BoolFlag{Name: "w2f,f", Usage: `-生成到文件`},
					},
				},
			},
		})
}
