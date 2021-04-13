package app

const tmplConfGo = `package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/conf/vars/db"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	
	//设置配置参数
	hydra.Conf.Web("8089").Header(header.WithCrossDomain()).Static(static.WithAssetsPath("./static"))
	hydra.Conf.Vars().DB().MySQL("db", "root", "rTo0CesHi2018Qx", "192.168.0.36:3306", "sms_test", db.WithConnect(20, 10, 600))
  //hydra.Conf.Vars().DB().Oracle("db", "sso", "123456", "orcl136", db.WithConnect(20, 10, 600))
	
	//启动时参数配置检查
	App.OnStarting(func(appConf app.IAPPConf) error {

		if _, err := hydra.C.DB().GetDB(); err != nil {
			return fmt.Errorf("db数据库配置错误,err:%v", err)
		}

		return nil
	})
}

`

const tmplConfSSoGo = `package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/sso/sso"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {

	//设置配置参数
	hydra.Conf.Web("8089").Header(header.WithCrossDomain()).
		Jwt(jwt.WithMode("HS512"),
		jwt.WithSecret("be29d433784c06dc8815ba7e8bcf35f3"),
		jwt.WithName("__jwt__"),
		jwt.WithExpireAt(36000),
		jwt.WithHeader(),
		jwt.WithExcludes("/sso/login/verify")).
		Sub("app", {###}
		{	
		"sso_secret":"B128F779D5741E701923346F7FA9F95C",
		"sso_api_host":"http://ssov4.100bm1.com:6689",
		"sso_ident":"sso_ident"
		}{###})

	hydra.Conf.Vars().DB().MySQL("db", "root", "rTo0CesHi2018Qx", "192.168.0.36:3306", "sms_test", db.WithConnect(20, 10, 600))

	//每个请求执行前执行
	App.OnHandleExecuting(func(ctx hydra.IContext) (rt interface{}) {
		if err := sso.CheckAndSetMember(ctx); err != nil {
			return err
		}
		return nil
	})

	//启动时参数配置检查
	App.OnStarting(func(c app.IAPPConf) error {
		conf := make(map[string]string)
		if _, err := c.GetServerConf().GetSubObject("app", conf); err != nil {
			return err
		}

		//检查db配置是否正确
		if _, err := hydra.C.DB().GetDB(); err != nil {
			return fmt.Errorf("db数据库配置错误,err:%v", err)
		}

		if err := sso.Config(conf["sso_api_host"], conf["sso_ident"], conf["sso_secret"], sso.WithAuthorityIgnore("/sso/**")); err != nil {
			return err
		}

		return nil
	})

}
`

const SnippetTmplConfGo = `package {{if (hasPrefix .GOMOD .ProjectPath )}}main{{else}}{{.ProjectPath|fileBasePath}}{{end}}

import (
	"github.com/micro-plat/hydra"
	{{- range $i,$v:=.Confs|importPath }}
	{{or $v.PkGAlias ""}}"{{$i}}"
	{{- end}}
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	hydra.OnReady(func() {
	{{- range $i,$v:=.Confs }}
		hydra.S.Web("/{{$v.Name|rmhd|rpath}}", {{$v.PKG}}.New{{$v.Name|rmhd|varName}}Handler())
	{{- end}}
	})
}

`
