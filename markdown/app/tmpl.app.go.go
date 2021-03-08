package app

const tmplAppGo = `package main
import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/app"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	//设置配置参数
	install()

	//启动时参数配置检查
	App.OnStarting(func(appConf app.IAPPConf) error {

		if _, err := hydra.C.DB().GetDB(); err != nil {
			return fmt.Errorf("db数据库配置错误,err:%v", err)
		}

		return nil
	})
}

`
