package main

import (
	"server"

	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/module"
)

func main() {

	app := mqant.CreateApp()
	app.OnConfigurationLoaded(func(app module.App) {

	})
	app.OnStartup(func(app module.App) {

	})
	app.Run(false, //只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		server.Module(),
	)

}
