package server

import (
	"github.com/urfave/cli"
)

func runServer() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		//构建服务
		s, err := newServer(c)
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
