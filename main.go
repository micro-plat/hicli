package main

import (
	"github.com/lib4dev/cli"
	_ "github.com/micro-plat/hicli/clones"
	_ "github.com/micro-plat/hicli/email"
	_ "github.com/micro-plat/hicli/markdown"
	"github.com/micro-plat/lib4go/logger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/micro-plat/hicli/pulls"
	_ "github.com/micro-plat/hicli/resets"
	_ "github.com/micro-plat/hicli/server"
	_ "github.com/micro-plat/hicli/update"
)

func main() {
	logger.Pause()
	var app = cli.New(cli.WithVersion("0.1.1"))
	app.Start()
}
